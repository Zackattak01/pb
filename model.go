package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + " " + i.desc }

type ListKeyMap struct {
    selectDirectory key.Binding
    goBack key.Binding
}

func newListKeyMap() *ListKeyMap {
    return &ListKeyMap{
       selectDirectory: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select directory")), 
       goBack: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "go back")), 
    }
}

type mode int 

const (
    OpenAsDirectory mode = 0
    OpenAsProject mode = 1
)

type model struct {
	list list.Model
    currentPath string
    keys *ListKeyMap
    mode mode
    depth int
    settings Settings
    options Options
}

func NewModel(settings Settings, options Options) model {
    keys := newListKeyMap()
    list := NewDirectoryList(settings, keys)
    return model{
        list: list,
        currentPath: "",
        keys: keys,
        mode: OpenAsDirectory,
        depth: 0,
        settings: settings,
        options: options,
    }
}

func (mod model) Init() (tea.Model, tea.Cmd) {
    argCount := len(mod.options.PositionalArguments)
    if mod.options.CreateTempProject {
        title := "tmp"
        if argCount >= 1 {
            title = mod.options.PositionalArguments[0]
        }

        path := "."
        if  argCount >= 2 {
            path = mod.options.PositionalArguments[1]
        }

        return mod, openProject(title, path, mod.settings.ProjectOpenCommand)
    } else if argCount == 1 {
        query := mod.options.PositionalArguments[0]

        for _, val := range mod.list.Items() {
            item, ok := val.(item)
            if ok && strings.EqualFold(item.title, query) {
                return mod, openProject(item.title, item.desc, mod.settings.ProjectOpenCommand)
            }
        }

        mod.list.SetFilterText(query)
    }
    
    return mod, nil
}

func (mod model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case projectClosedMsg:
        if mod.options.QuitOnProjectExit {
            return mod, tea.Quit
        }
	case tea.KeyMsg:
        if mod.list.FilterState() == list.Filtering {
            break
        }

        switch { 
        case key.Matches(msg, mod.keys.selectDirectory):
            item, ok := mod.list.SelectedItem().(item)
            if ok {
                if mod.mode == OpenAsDirectory {
                    // we store the absolute path of the item in the description
                    mod.currentPath = item.desc
                    mod.list.SetItems(NewDirectoryListItems(Source {TraversalDepth: 0, Path: mod.currentPath}))
                    mod.depth++
                } else if mod.mode == OpenAsProject {
                    // we store the absolute path of the item in the description
                    mod.list.ResetFilter()
                    return mod, openProject(item.title, item.desc, mod.settings.ProjectOpenCommand)
                }
            }

        case key.Matches(msg, mod.keys.goBack):
            mod.list.ResetFilter()
            if mod.depth <= 0 {
                return mod, nil
            }

            mod.depth--
            if mod.depth == 0 {
                mod.currentPath = ""
                mod.list.SetItems(NewDirectoryListItems(mod.settings.Sources...))
            } else {
                mod.currentPath = filepath.Join(mod.currentPath, "..")
                mod.list.SetItems(NewDirectoryListItems(Source {TraversalDepth: 0, Path: mod.currentPath}))
            }
        }
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		mod.list.SetSize(msg.Width-h, msg.Height-v)
	}

    if mod.depth == mod.settings.DefaultOpenDepth {
        mod.mode = OpenAsProject
    } else {
        mod.mode = OpenAsDirectory
    }

	var cmd tea.Cmd
	mod.list, cmd = mod.list.Update(msg)
	return mod, cmd
}

func (mod model) View() string {
    originalWidth := mod.list.Width()
    
    largestItemWidth := 0
    
    for _, i := range mod.list.Items() {
        item := i.(item)
        
        if len(item.desc) > largestItemWidth {
            largestItemWidth = len(item.desc)
            
        }
    }
    largestItemWidth += 4

    mod.list.SetWidth(largestItemWidth)
    
    marginLen := (originalWidth - largestItemWidth) / 2
    style := lipgloss.NewStyle().Width(largestItemWidth).Margin(0, marginLen)

    mod.list.Styles.TitleBar = mod.list.Styles.TitleBar.Width(largestItemWidth).AlignHorizontal(lipgloss.Center)
    mod.list.Styles.StatusBar = mod.list.Styles.StatusBar.Width(largestItemWidth).AlignHorizontal(lipgloss.Center)
     
    return style.Render(mod.list.View())
}

type projectClosedMsg struct{ err error }

func openProject(name, path, defaultCommand string) tea.Cmd {
    command := defaultCommand
    variables := strings.NewReplacer("$projectName", name, "$projectPath", path)

    projectConfig, err := LoadProjectConfig(path)
    if err == nil {
        variables = strings.NewReplacer("$projectName", name, "$projectPath", path, "$defaultCommand", variables.Replace(command))
        command = projectConfig.ProjectOpenCommand
    }

    parsedCommand := parseCommand(command, variables)
    c := exec.Command(parsedCommand[0], parsedCommand[1:]...)
    c.Env = os.Environ()
    c.Env = append(c.Env, "PB_NAME=" + name)
    c.Env = append(c.Env, "PB_PATH=" + path)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return projectClosedMsg{err}
	})
}

func parseCommand(command string, variables *strings.Replacer) []string {
    command = variables.Replace(command)
    seperatedCommand := make([]string, 0, 5)
    lastStrEnd := 0
    openQuote := rune(0)
    processQuote := false
    
    for i, c := range command {
        if c == '\'' || c == '"' {
            if openQuote == c {
                openQuote = rune(0)
                processQuote = true
            } else {
                openQuote = c
            }
        }

        if openQuote == rune(0) && c == ' ' {
            var arg string
            if processQuote {
                arg = command[lastStrEnd+1:i-1]
                processQuote = false
            } else {
                arg = command[lastStrEnd:i]
            }

            seperatedCommand = append(seperatedCommand, arg)
            lastStrEnd = i + 1
        }
    }

    if lastStrEnd < len(command) {
        var arg string
        if processQuote {
            arg = command[lastStrEnd+1:len(command)-1]
        } else {
            arg = command[lastStrEnd:]
        }

        seperatedCommand = append(seperatedCommand, arg)
    }
    return seperatedCommand
}

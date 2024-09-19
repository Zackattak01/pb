package main

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
}

func NewModel(settings Settings) model {
    keys := newListKeyMap()
    list := NewDirectoryList(settings, keys)
    return model{
        list: list,
        currentPath: "",
        keys: keys,
        mode: OpenAsDirectory,
        depth: 0,
        settings: settings,
    }
}

func (m model) Init() tea.Cmd {
	return nil
}

func (mod model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
                    mod.list.SetItems(NewDirectoryListItems(0, mod.currentPath))
                    mod.depth++
                } else if mod.mode == OpenAsProject {
                    // we store the absolute path of the item in the description
                    mod.list.FilterInput.SetValue("")
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
                mod.list.SetItems(NewDirectoryListItems(mod.settings.SourceTraversalDepth, mod.settings.Sources...))
            } else {
                mod.currentPath = filepath.Join(mod.currentPath, "..")
                mod.list.SetItems(NewDirectoryListItems(0, mod.currentPath))
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
	return docStyle.Render(mod.list.View())
}

type projectClosedMsg struct{ err error }

func openProject(name, path, command string) tea.Cmd {
    parsedCommand := parseCommand(command, strings.NewReplacer("$projectName", name, "$projectPath", path))

    c := exec.Command(parsedCommand[0], parsedCommand[1:]...)

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

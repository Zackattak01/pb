package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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
        switch { 
        case key.Matches(msg, mod.keys.selectDirectory):
            item, ok := mod.list.SelectedItem().(item)
            if ok {
                if mod.mode == OpenAsDirectory {
                    // we store the absolute path of the item in the description
                    mod.currentPath = item.desc
                    mod.list.SetItems(NewDirectoryListItems(mod.currentPath))
                    mod.depth++
                } else if mod.mode == OpenAsProject {
                    // we store the absolute path of the item in the description
                    return mod, openEditor(item.desc)
                }
            }

        case key.Matches(msg, mod.keys.goBack):
            if mod.depth <= 0 {
                return mod, nil
            }

            mod.depth--
            if mod.depth == 0 {
                mod.currentPath = ""
                mod.list.SetItems(NewDirectoryListItems(mod.settings.Sources...))
            } else {
                mod.currentPath = filepath.Join(mod.currentPath, "..")
                mod.list.SetItems(NewDirectoryListItems(mod.currentPath))
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

type editorFinishedMsg struct{ err error }

func openEditor(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	c := exec.Command(editor, "-c", fmt.Sprintf("cd %s", path)) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

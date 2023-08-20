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

type Mode int 

const (
    List Mode = 0
    OpenAsProject Mode = 1
)

type model struct {
	list list.Model
    currentPath string
    keys *ListKeyMap
    mode Mode
    depth int
}

func NewModel(path string) model {
    keys := newListKeyMap()
    list := NewDirectoryList(path, keys)
    return model{
        list: list,
        currentPath: path,
        keys: keys,
        mode: List,
        depth: 0,
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
                if mod.mode == List {
                    mod.currentPath = filepath.Join(mod.currentPath, item.title) 
                    mod.list.SetItems(NewDirectoryListItems(mod.currentPath))
                    mod.depth++
                } else if mod.mode == OpenAsProject {
                    return mod, openEditor(filepath.Join(mod.currentPath, item.title))
                }
            }

        case key.Matches(msg, mod.keys.goBack):
            mod.currentPath = filepath.Join(mod.currentPath, "..")
            mod.list.SetItems(NewDirectoryListItems(mod.currentPath))
            mod.depth--
        }
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		mod.list.SetSize(msg.Width-h, msg.Height-v)
	}

    if mod.depth == 1 {
        mod.mode = OpenAsProject
    } else {
        mod.mode = List
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

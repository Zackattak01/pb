package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func NewDirectoryList(path string, keys *ListKeyMap) list.Model {
    list := list.New(NewDirectoryListItems(path), createDelegate(), 0, 0)
    list.Title = "Projects"
    list.InfiniteScrolling = true
    list.AdditionalShortHelpKeys = func() []key.Binding {
        return []key.Binding{
            keys.selectDirectory,
            keys.goBack,
        }
    }

    list.SetShowPagination(false)

    return list 
}

func NewDirectoryListItems(path string) []list.Item {
    directories := getDirectories(path)
    items := make([]list.Item, len(directories))

    for i, dir := range directories { 
        items[i] = item{title: dir.Name(), desc: ""}
    }

    return items
}

func getDirectories(path string) []os.DirEntry {
    contents, err := os.ReadDir(path)
    if err != nil {
        log.Fatalf("Could not read files from path: %s", path)
    }

    directories := make([]os.DirEntry, 0, len(contents))
    for _, file := range contents {
        if file.IsDir() {
             directories = append(directories, file)
        }
    }

    return directories
}

func createDelegate() list.DefaultDelegate {
    delegate := list.NewDefaultDelegate()
    delegate.ShowDescription = false

    return delegate
}

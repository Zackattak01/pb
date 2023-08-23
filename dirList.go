package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func NewDirectoryList(settings Settings, keys *ListKeyMap) list.Model {
    list := list.New(NewDirectoryListItems(settings.Sources...), createDelegate(settings), 0, 0)
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

func NewDirectoryListItems(paths ...string) []list.Item {
    directories := getDirectories(paths...)
    items := make([]list.Item, len(directories))

    for i, dir := range directories { 
        items[i] = item{title: dir.name, desc: dir.path}
    }

    return items
}
 
type directory struct {
    name string
    path string
}


func getDirectories(paths ...string) []directory {
    // create empty slice with default cap of 10
    directories := make([]directory, 0, 10)

    for _, path := range paths {
        contents, err := os.ReadDir(path)
        if err != nil {
            log.Fatalf("Could not read files from path: %s", paths)
        }

        for _, file := range contents {
            if file.IsDir() {
                 directories = append(directories, directory{name: file.Name(), path: filepath.Join(path, file.Name())})
            }
        }
    }

    return directories
}

func createDelegate(settings Settings) list.DefaultDelegate {
    delegate := list.NewDefaultDelegate()
    delegate.ShowDescription = settings.DisplayAbsolutePath

    return delegate
}

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func NewDirectoryList(settings Settings, keys *ListKeyMap) list.Model {
    list := list.New(NewDirectoryListItems(settings.SourceTraversalDepth, settings.Sources...), createDelegate(settings), 0, 0)
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

func NewDirectoryListItems(traversalDepth int, paths ...string) []list.Item {
    directories := getDirectories(traversalDepth, paths...)
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


func getDirectories(traversalDepth int, paths ...string) []directory {
    directories := make([]directory, 0, 10)
    traversePaths := make([]string, len(paths))
    copy(traversePaths, paths)

    for i:= 0;  i < traversalDepth + 1;  i++ {
        newPaths := make([]string, 0, len(traversePaths)) 
        for _, path := range traversePaths {
            contents, err := os.ReadDir(path)
            if err != nil {
                log.Fatalf("Could not read files from path: %s", paths)
            }
            
            for _, file := range contents {
                if file.IsDir() {
                    filePath := filepath.Join(path, file.Name())
                   if i == traversalDepth {
                       directories = append(directories, directory{name: file.Name(), path: filePath})
                   } else {
                       newPaths = append(newPaths, filePath)
                   }
                }
            }
        }

        traversePaths = newPaths
    }

    return directories
}

func createDelegate(settings Settings) list.DefaultDelegate {
    delegate := list.NewDefaultDelegate()
    delegate.ShowDescription = settings.DisplayAbsolutePath

    return delegate
}

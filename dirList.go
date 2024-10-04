package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
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

func NewDirectoryListItems(sources ...Source) []list.Item {
    directories := getDirectories(sources...)
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


func getDirectories(sources ...Source) []directory {
    directories := make([]directory, 0, 10)

    // yes there are 4 for loops here
    for _, source := range sources {
        sourcePaths := make([]string, 0, 20)
        sourcePaths = append(sourcePaths, source.Path)

        for i:= 0;  i < source.TraversalDepth + 1;  i++ {
            newPaths := make([]string, 0, len(sourcePaths)) 
            for _, path := range sourcePaths {
                contents, err := os.ReadDir(path)
                if err != nil {
                    log.Fatalf("Could not read files from path: %s", sources)
                }

                for _, file := range contents {
                    if file.IsDir() {
                        filePath := filepath.Join(path, file.Name())
                        if i == source.TraversalDepth {
                            directories = append(directories, directory{name: file.Name(), path: filePath})
                        } else {
                            newPaths = append(newPaths, filePath)
                        }
                    }
                }
            }

            sourcePaths = newPaths
        }

    }

    return directories
}

func createDelegate(settings Settings) list.DefaultDelegate {
    delegate := list.NewDefaultDelegate()
    delegate.ShowDescription = settings.DisplayAbsolutePath

    return delegate
}

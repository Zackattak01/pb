package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const SettingsPath = "pb/settings.json"
const ProjectConfigFile = ".pb.json"

type Source struct {
    Path string
    TraversalDepth int
}

type Settings struct {
    Sources []Source
    ProjectOpenCommand string
    DefaultOpenDepth int
    DisplayAbsolutePath bool
}

type ProjectConfig struct {
    ProjectOpenCommand string
}

type Options struct {
    PositionalArguments []string
    CreateTempProject bool
    QuitOnProjectExit bool
}

func newOptions() Options {
    return Options{CreateTempProject: false, QuitOnProjectExit: false, PositionalArguments: make([]string, 0)}
}

func LoadSettings() (*Settings, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return nil, err
    }

    settingsPath := filepath.Join(configDir, SettingsPath)
    contents, err := os.ReadFile(settingsPath)
    if err != nil {
        return nil, err
    }

    var settings Settings
    err = json.Unmarshal(contents, &settings)
    if err != nil {
        return nil, err
    }

    return &settings, nil 
}

func LoadProjectConfig(projectDir string) (*ProjectConfig, error) {
    configPath := filepath.Join(projectDir, ProjectConfigFile)
    contents, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config ProjectConfig
    err = json.Unmarshal(contents, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil 
}

func ParseOptions(args []string) (*Options, error) {
    options := newOptions()
    for _, arg := range args {
        if arg == "-t" || arg == "--temp" {
            options.CreateTempProject = true
        } else if arg == "-q" || arg == "-q" {
            options.QuitOnProjectExit = true;
        } else {
            options.PositionalArguments = append(options.PositionalArguments, arg)
        }
    }

    return &options, nil
}

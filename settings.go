package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const SettingsPath = "pb/settings.json"
const ProjectConfigFile = ".pb.json"

const indent = "    "

var TempPath = filepath.Join(os.TempDir(), "pb")

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

    // add an additional source so tmp projects can be recovered
    settings.Sources = append(settings.Sources, Source{Path: TempPath, TraversalDepth: 0})
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

func WriteProjectConfig(dir string, config ProjectConfig) error {
    jsonString, err := json.MarshalIndent(config, "", indent)
    if err != nil {
        return err
    }

    path := filepath.Join(dir, ProjectConfigFile)
    err = os.WriteFile(path, jsonString, 0700)
    if err != nil {
        return err
    }
    
    return nil
}

func ParseOptions(args []string) (*Options, error) {
    options := newOptions()
    for _, arg := range args {
        if arg == "-t" || arg == "--temp" {
            options.CreateTempProject = true
        } else if arg == "-q" || arg == "--quit" || arg == "--quit-on-close" {
            options.QuitOnProjectExit = true;
        } else {
            options.PositionalArguments = append(options.PositionalArguments, arg)
        }
    }

    return &options, nil
}

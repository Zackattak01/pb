package main

import (
	"encoding/json"
	"os"
	"path"
)

const SettingsPath = "pb/settings.json"

type Settings struct {
    Sources []string
    SourceTraversalDepth int
    ProjectOpenCommand string
    DefaultOpenDepth int
    DisplayAbsolutePath bool
}

func LoadSettings() (*Settings, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return nil, err
    }

    settingsPath := path.Join(configDir, SettingsPath)
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

package main

import (
	"encoding/json"
	"log"
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

func LoadSettings() Settings {
    configDir, err := os.UserConfigDir()

    if err != nil {
        log.Fatal("Could not get config dir")
    }

    settingsPath := path.Join(configDir, SettingsPath)
    contents, err := os.ReadFile(settingsPath)

    if err != nil {
        log.Fatal("Error reading settings file")
    }

    var settings Settings
    err = json.Unmarshal(contents, &settings)
    if err != nil {
        log.Fatal("Error reading settings file")
    }

    return settings 
}

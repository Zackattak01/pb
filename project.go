package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type ProjectClosedMsg struct{ err error }

func OpenTempProject(name, path, defaultCommand string) tea.Cmd {
    path, _ = filepath.Abs(path)
    variables := getVariableReplacer(name, path);
    command := substituteVariables(defaultCommand, variables)

    // write a project config to a temp dir so that PB can recover the temp sessions
    tempDir, err := getTempProjectDir(name)
    if err == nil {
        WriteProjectConfig(tempDir, ProjectConfig{ProjectOpenCommand: command})
    }

    return execProject(name, path, command)
}

func getTempProjectDir(name string) (string, error) {
    tmpFolder := filepath.Join(os.TempDir(), "pb")
    if _, err := os.Stat(tmpFolder); err != nil {
        if err := os.Mkdir(tmpFolder, 0700); err != nil {
            return "", err
        }
    }

    tmpProjFolder := filepath.Join(tmpFolder, name)
    if err := os.Mkdir(tmpProjFolder, 0700); err != nil {
        return "", err
    }

    return tmpProjFolder, nil
}

func OpenProject(name, path, defaultCommand string) tea.Cmd {
    command := defaultCommand
    variables := getVariableReplacer(name, path);

    projectConfig, err := LoadProjectConfig(path)
    if err == nil {
        command = strings.ReplaceAll(projectConfig.ProjectOpenCommand, "$PB_CMD", command)
    }

    return execProject(name, path, substituteVariables(command, variables))
}

func execProject(name, path, command string) tea.Cmd {
    argumentizedCommand := argumentizeCommand(command)
    c := exec.Command(argumentizedCommand[0], argumentizedCommand[1:]...)
    c.Env = os.Environ()
    c.Env = append(c.Env, "PB_NAME=" + name)
    c.Env = append(c.Env, "PB_PATH=" + path)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return ProjectClosedMsg{err}
	})
}

func getVariableReplacer(name, path string) *strings.Replacer {
    return strings.NewReplacer("$PB_NAME", name, "$PB_PATH", path)
}

func substituteVariables(command string, variables *strings.Replacer) string {
    return variables.Replace(command)
}

func argumentizeCommand(command string) []string {
    argumentizedCommand := make([]string, 0, 5)
    lastStrEnd := 0
    openQuote := rune(0)
    processQuote := false
    
    for i, c := range command {
        if c == '\'' || c == '"' {
            if openQuote == c {
                openQuote = rune(0)
                processQuote = true
            } else {
                openQuote = c
            }
        }

        if openQuote == rune(0) && c == ' ' {
            var arg string
            if processQuote {
                arg = command[lastStrEnd+1:i-1]
                processQuote = false
            } else {
                arg = command[lastStrEnd:i]
            }

            argumentizedCommand = append(argumentizedCommand, arg)
            lastStrEnd = i + 1
        }
    }

    if lastStrEnd < len(command) {
        var arg string
        if processQuote {
            arg = command[lastStrEnd+1:len(command)-1]
        } else {
            arg = command[lastStrEnd:]
        }

        argumentizedCommand = append(argumentizedCommand, arg)
    }
    return argumentizedCommand
}

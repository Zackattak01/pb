package main

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type ProjectClosedMsg struct{ err error }

func OpenProject(name, path, defaultCommand string) tea.Cmd {
    command := defaultCommand
    variables := strings.NewReplacer("$projectName", name, "$projectPath", path)

    projectConfig, err := LoadProjectConfig(path)
    if err == nil {
        variables = strings.NewReplacer("$projectName", name, "$projectPath", path, "$defaultCommand", variables.Replace(command))
        command = projectConfig.ProjectOpenCommand
    }

    parsedCommand := parseCommand(command, variables)
    c := exec.Command(parsedCommand[0], parsedCommand[1:]...)
    c.Env = os.Environ()
    c.Env = append(c.Env, "PB_NAME=" + name)
    c.Env = append(c.Env, "PB_PATH=" + path)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return ProjectClosedMsg{err}
	})
}

func parseCommand(command string, variables *strings.Replacer) []string {
    command = variables.Replace(command)
    seperatedCommand := make([]string, 0, 5)
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

            seperatedCommand = append(seperatedCommand, arg)
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

        seperatedCommand = append(seperatedCommand, arg)
    }
    return seperatedCommand
}

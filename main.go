package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle()


func main() {
    settings, err := LoadSettings()
    if err != nil {
        logFatalError("Error reading settings:", err)
    }

    mod := NewModel(*settings)

	program := tea.NewProgram(mod, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
        logFatalError("Error running program:", err)
	}
}


func logFatalError(msg string, err error) {
    fmt.Println(msg, err)
    os.Exit(1)
}

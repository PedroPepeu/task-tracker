// main.go

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"task-tracker/storage"
	"task-tracker/ui"
)

func main() {
	tasks := storage.LoadTask("data.json")

	initialState := ui.InitialModel()
	initialState.Choices = tasks
	initialState.FilteredChoices = tasks

	p := tea.NewProgram(initialState)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

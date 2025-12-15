package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/huh"
)

type task_status int

const (
	todo task_status = iota
	inprogress
	done
)

const fileName = "data.json"

type task struct {
	id          int         `json:"id"`
	description string      `json:"description"`
	status      task_status `json:"status"`
	createdAt   time.Time   `json:"created"`
	updatedAt   time.Time   `json:"updated"`
}

type model struct {
	choices  []task
	cursor   int
	selected map[int]struct{}
}

func saveTask(filename string, data []task) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

func loadTask(filename string) []task {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("file not found or error reading:", err)
		return []task{}
	}

	var tasks []task

	err = json.Unmarshal(file, &tasks)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	return tasks
}

form := huh.NewForm(
	huh.NewGroup(
		huh.NewText().
			Title("Title of the task").
			CharLimit(40).
			Value(&title),

		huh.NewText().
			Title("Description of the task").
			CharLimit(200).
			Value($description),

		huh.NewConfirm().
			Title("Wanna create this task?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm)
	),
)

func initialModel() model {
	return model{
		choices: []task{}, // examples

		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What tasks do you plan to solve today?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

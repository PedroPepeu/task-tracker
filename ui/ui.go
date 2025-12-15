// ui.go
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"task-tracker/model"
)

const (
	StateList = iota
	StateForm
)

type MainModel struct {
	State    int
	Choices  []model.Task
	Cursor   int
	Selected map[int]struct{}

	form *huh.Form
}

func InitialModel() MainModel {
	return MainModel{
		State:    StateList,
		Choices:  []model.Task{}, // examples
		Selected: make(map[int]struct{}),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func NewTaskForm() *huh.NewForm {
	var title string
	var desc string
	var confirm bool

	return huh.newForm(
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
	).WithTheme(huh.ThemeBase())
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		case "enter", " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m MainModel) View() string {
	s := "What tasks do you plan to solve today?\n\n"

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.Selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Description)
	}

	s += "\nPress q to quit.\n"

	return s
}

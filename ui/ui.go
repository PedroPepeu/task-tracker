// ui.go
package ui

import (
	"fmt"
	"time"

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

	formTitle string
	formDesc  string
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

func NewTaskForm(title *string, desc *string) *huh.Form {
	var confirm bool

	return huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Title of the task").
				CharLimit(40).
				Value(title),

			huh.NewText().
				Title("Description of the task").
				CharLimit(200).
				Value(desc),

			huh.NewConfirm().
				Title("Wanna create this task?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeBase())
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// FORM MODE
	if m.State == StateForm {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f

			if m.form.State == huh.StateCompleted {
				newTask := model.Task{
					ID:          len(m.Choices) + 1,
					Title:       m.formTitle,
					Description: m.formTitle,
					CreatedAt:   time.Now(),
					Status:      model.Todo,
				}

				m.Choices = append(m.Choices, newTask)

				m.State = StateList
				m.formTitle = ""
				m.formDesc = ""
			}
		}
		return m, cmd
	}

	// LIST MODE

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

		case "n":
			m.State = StateForm
			m.form = NewTaskForm(&m.formTitle, &m.formDesc)
			return m, m.form.Init()
		}
	}

	return m, nil
}

func (m MainModel) View() string {
	if m.State == StateForm {
		return m.form.View()
	}

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

		s += fmt.Sprintf("%s [%s] %s - %s\n", cursor, checked, choice.Title, choice.Description)
	}

	s += "\n---------------------------"
	s += "\n[ N ] New Task  |  [ Q ] Quit\n"

	return s
}

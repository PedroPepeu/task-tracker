// ui.go
package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"task-tracker/model"
	"task-tracker/storage"
)

const (
	StateList = iota
	StateForm
)

type FormData struct {
	Title   string
	Desc    string
	Confirm bool
}

type MainModel struct {
	State    int
	Choices  []model.Task
	Cursor   int
	Selected map[int]struct{}

	form *huh.Form

	formData *FormData
}

func InitialModel() MainModel {
	return MainModel{
		State:    StateList,
		Choices:  []model.Task{}, // examples
		Selected: make(map[int]struct{}),

		formData: &FormData{},
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func NewTaskForm(data *FormData) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title of the task").
				CharLimit(40).
				Value(&data.Title),

			huh.NewText().
				Title("Description of the task").
				CharLimit(200).
				Value(&data.Desc),

			huh.NewConfirm().
				Title("Wanna create this task?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&data.Confirm),
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
				if m.formData.Confirm && m.formData.Title != "" {
					newTask := model.Task{
						ID:          len(m.Choices) + 1,
						Title:       m.formData.Title,
						Description: m.formData.Desc,
						CreatedAt:   time.Now(),
						Status:      model.Todo,
					}

					m.Choices = append(m.Choices, newTask)
					storage.SaveTask("data.json", m.Choices)
				} else {
					fmt.Println("DEBUG ERROR: Not saving because Title is empty or Confirm is No")
				}

				m.State = StateList
				m.formData.Title = ""
				m.formData.Desc = ""
				m.formData.Confirm = false
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
			m.form = NewTaskForm(m.formData)
			return m, m.form.Init()
		}
	}

	return m, nil
}

func (m MainModel) View() string {
	if m.State == StateForm {
		return m.form.View()
	}

	s := fmt.Sprintf("Debug: I have %d tasks loaded.\n", len(m.Choices))
	s += "What tasks do you plan to solve today?\n\n"

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
	s += "\n[ N ] New Task  |  [ E ] Edit Task  |  [ Q ] Quit\n"

	return s
}

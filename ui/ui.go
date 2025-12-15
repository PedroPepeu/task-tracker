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

type MainModel struct {
	State    int
	Choices  []model.Task
	Cursor   int
	Selected map[int]struct{}

	form     *huh.Form
	formData *FormData

	indexToEdit int
}

func InitialModel() MainModel {
	return MainModel{
		State:       StateList,
		Choices:     []model.Task{},
		Selected:    make(map[int]struct{}),
		formData:    &FormData{},
		indexToEdit: -1,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// FORM MODE
	if m.State == StateForm {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f

			if m.form.State == huh.StateCompleted {
				if m.formData.Confirm && m.formData.Title != "" {
					if m.indexToEdit == -1 {
						newTask := model.Task{
							ID:          len(m.Choices) + 1,
							Title:       m.formData.Title,
							Description: m.formData.Desc,
							CreatedAt:   time.Now(),
							Status:      model.Todo,
						}

						m.Choices = append(m.Choices, newTask)
					} else {
						m.Choices[m.indexToEdit].Title = m.formData.Title
						m.Choices[m.indexToEdit].Description = m.formData.Desc
						m.Choices[m.indexToEdit].UpdatedAt = time.Now()
					}

					storage.SaveTask("data.json", m.Choices)
				}

				m.State = StateList
				m.formData.Title = ""
				m.formData.Desc = ""
				m.formData.Confirm = false
				m.indexToEdit = -1
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
			if len(m.Choices) > 0 {
				currentStatus := m.Choices[m.Cursor].Status

				if currentStatus == model.Done {
					m.Choices[m.Cursor].Status = model.Todo
				} else {
					m.Choices[m.Cursor].Status++
				}

				storage.SaveTask("data.json", m.Choices)
			}

		case "n":
			m.indexToEdit = -1

			m.formData.Title = ""
			m.formData.Desc = ""
			m.formData.Confirm = false

			m.State = StateForm
			m.form = NewTaskForm(m.formData)
			return m, m.form.Init()

		case "d":
			if len(m.Choices) > 0 {
				m.Choices = append(m.Choices[:m.Cursor], m.Choices[m.Cursor+1:]...)

				storage.SaveTask("data.json", m.Choices)

				if m.Cursor >= len(m.Choices) && m.Cursor > 0 {
					m.Cursor--
				}

				m.Selected = make(map[int]struct{})
			}

		case "e":
			if len(m.Choices) > 0 {
				m.indexToEdit = m.Cursor

				taskToEdit := m.Choices[m.Cursor]

				m.formData.Title = taskToEdit.Title
				m.formData.Desc = taskToEdit.Description
				m.formData.Confirm = false

				m.State = StateForm
				m.form = NewTaskForm(m.formData)
				return m, m.form.Init()
			}
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

		icon := " "
		switch choice.Status {
		case model.Todo:
			icon = "To Do"

		case model.Inprogress:
			icon = "In Progress"

		case model.Done:
			icon = "Done"
		}

		s += fmt.Sprintf("%s [%s] %s - %s\n", cursor, icon, choice.Title, choice.Description)
	}

	s += "\n-----------------------------------------------------------------------"
	s += "\n[ N ] New Task  |  [ E ] Edit Task  |  [ D ] Delete Task  |  [ Q ] Quit\n"

	return s
}

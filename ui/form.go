// form.go
package ui

import (
	"github.com/charmbracelet/huh"
)

type FormData struct {
	Title   string
	Desc    string
	Confirm bool
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

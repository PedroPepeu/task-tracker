// utils.go
package ui

import (
	"task-tracker/model"
)

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func applyFilter(choices []model.Task, state int) []model.Task {
	if state == -1 {
		return choices
	}

	rule := func(t model.Task) bool {
		return int(t.Status) == state
	}

	return filter(choices, rule)
}

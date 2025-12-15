// task.go
package model

import (
	"time"
)

type TaskStatus int

const (
	Todo TaskStatus = iota
	Inprogress
	Done
)

const fileName = "data.json"

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created"`
	UpdatedAt   time.Time  `json:"updated"`
}

func (t Task) String() string {
	return t.Description
}

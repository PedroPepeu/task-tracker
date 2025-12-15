// task.go
package model

import (
	"fmt"
	"time"
)

type TaskStatus int

const (
	Todo TaskStatus = iota
	Inprogress
	Done
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created"`
	UpdatedAt   time.Time  `json:"updated"`
}

func (t Task) String() string {
	return fmt.Sprintf("%s: %s", t.Title, t.Description)
}

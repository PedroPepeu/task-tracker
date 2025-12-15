// task_test.go
package model

import (
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	newTask := Task{
		ID:          1,
		Title:       "Test Task",
		Description: "This is a test",
		Status:      Todo,
		CreatedAt:   time.Now(),
	}

	if newTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", newTask.Title)
	}

	if newTask.Status != Todo {
		t.Errorf("Expected status Todo (0), got %d", newTask.Status)
	}
}

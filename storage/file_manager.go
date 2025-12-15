// file_manager.go
package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"task-tracker/model"
)

func SaveTask(filename string, data []model.Task) {
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

func LoadTask(filename string) []model.Task {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("file not found or error reading:", err)
		return []model.Task{}
	}

	var tasks []model.Task

	err = json.Unmarshal(file, &tasks)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	return tasks
}

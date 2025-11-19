package flaghandler

import (
	"fmt"

	filemanager "github.com/NikitaTumanov/terminalTaskTracker/internal/file_manager"
	"github.com/NikitaTumanov/terminalTaskTracker/internal/models"
)

type Storage struct {
	Tasks []models.Task
}

func New() (*Storage, error) {
	tasks, err := filemanager.GetAllTasks()
	if err != nil {
		return &Storage{}, fmt.Errorf("filemanager.GetAllTasks: %w", err)
	}

	return &Storage{
		Tasks: tasks,
	}, nil
}

func (s Storage) Handle() error {
	fmt.Println("321")
	return nil
}

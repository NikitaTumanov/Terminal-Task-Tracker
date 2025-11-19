package flaghandler

import (
	"fmt"
	"time"

	filemanager "github.com/NikitaTumanov/terminalTaskTracker/internal/file_manager"
	"github.com/NikitaTumanov/terminalTaskTracker/internal/models"
)

type storage struct {
	tasks []models.Task
}

func New() (*storage, error) {
	tasks, err := filemanager.GetAllTasks()
	if err != nil {
		return &storage{}, fmt.Errorf("filemanager.GetAllTasks: %w", err)
	}

	return &storage{
		tasks: tasks,
	}, nil
}

func (s *storage) Update() {
	go func() {
		var err error
		for {
			time.Sleep(models.TimeOut)
			s.tasks, err = filemanager.GetAllTasks()
			if err != nil {
				fmt.Println("filemanager.GetAllTasks: ", err)
			}
		}
	}()
}

func (s storage) Handle() error {
	fmt.Println("321")
	return nil
}

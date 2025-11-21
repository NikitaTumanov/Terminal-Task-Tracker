package flaghandler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	filemanager "github.com/NikitaTumanov/terminalTaskTracker/internal/file_manager"
	"github.com/NikitaTumanov/terminalTaskTracker/internal/models"
)

type storage struct {
	tasks      []models.Task
	taskIndex  int
	command    string
	taskName   string
	taskStatus string
	helpFlag   bool
}

func New(taskIndex int, command, taskName, taskStatus string, helpFlag bool) (*storage, error) {
	tasks, err := filemanager.GetAllTasks()
	if err != nil {
		return &storage{}, fmt.Errorf("filemanager.GetAllTasks: %w", err)
	}

	return &storage{
		tasks:      tasks,
		taskIndex:  taskIndex,
		command:    command,
		taskName:   taskName,
		taskStatus: taskStatus,
		helpFlag:   helpFlag,
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

func printHelp() {
	fmt.Println(`	Add Task: -c add --name="<Task name>"
	Update Task: -c update --index=<Task Index> --name="<New Task Name>" --status=<New Task Status>
		Task Statuses:
			0 - Не начато
			1 - В процессе
			2 - Выполнено
	Delete Task: -c delete --index=<Task Index>
	Update Task Status: -c updateStatus --index=<Task Index> --status=<New Task Status>
		Task Statuses:
			0 - Не начато
			1 - В процессе
			2 - Выполнено
	Show All Tasks: -c allTasks
	Show Done Tasks: -c doneTasks
	Show Not Done Tasks: -c notDoneTasks
	Show Tasks In Progress: -c inProgressTasks`)
}

func (s *storage) Handle() error {
	if s.helpFlag {
		printHelp()
		return nil
	}

	switch strings.ToLower(s.command) {
	case "add":
		if s.taskName == "" {
			return errors.New("отсутствует name в переданных аргументах")
		}
		result, err := filemanager.Add(&s.tasks, []string{s.taskName})
		if err != nil {
			return fmt.Errorf("filemanager.Add: %w", err)
		}
		fmt.Println(result)

	case "update":
		if s.taskIndex == -1 {
			return errors.New("отсутствует index в переданных аргументах")
		}
		if s.taskName == "" {
			return errors.New("отсутствует name в переданных аргументах")
		}
		if s.taskStatus == "" {
			return errors.New("отсутствует status в переданных аргументах")
		}
		result, err := filemanager.Update(&s.tasks, []string{strconv.Itoa(s.taskIndex), s.taskName, s.taskStatus})
		if err != nil {
			return fmt.Errorf("filemanager.Update: %w", err)
		}
		fmt.Println(result)

	case "delete":
		if s.taskIndex == -1 {
			return errors.New("отсутствует index в переданных аргументах")
		}
		result, err := filemanager.Delete(&s.tasks, []string{strconv.Itoa(s.taskIndex)})
		if err != nil {
			return fmt.Errorf("filemanager.Delete: %w", err)
		}
		fmt.Println(result)

	case "updatestatus":
		if s.taskIndex == -1 {
			return errors.New("отсутствует index в переданных аргументах")
		}
		if s.taskStatus == "" {
			return errors.New("отсутствует status в переданных аргументах")
		}
		result, err := filemanager.UpdateStatus(&s.tasks, []string{strconv.Itoa(s.taskIndex), s.taskStatus})
		if err != nil {
			return fmt.Errorf("filemanager.UpdateStatus: %w", err)
		}
		fmt.Println(result)

	case "alltasks":
		err := filemanager.AllTasks(&s.tasks)
		return err

	case "donetasks":
		err := filemanager.DoneTasks(&s.tasks)
		return err

	case "notdonetasks":
		err := filemanager.NotDoneTasks(&s.tasks)
		return err

	case "inprogresstasks":
		err := filemanager.InProgressTasks(&s.tasks)
		return err

	case "help":
	default:
		return errors.New("введена некорректная команда")
	}

	return nil
}

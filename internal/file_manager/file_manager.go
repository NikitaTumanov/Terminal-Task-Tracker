package filemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/models"
)

const (
	tasksPath = "tasks.json"
)

var (
	ErrInputElementsCount error = errors.New("incorrect number of arguments passed")
	errAtoi               error = errors.New("an invalid number was passed")
	errIncorrectStatus    error = errors.New("an incorrect task status was passed")
	ErrNameNotExists      error = errors.New("name is missing from the passed arguments")
	ErrIndexNotExists     error = errors.New("index is missing from the passed arguments")
	ErrStatusNotExists    error = errors.New("status is missing from the passed arguments")
	ErrInvalidCommand     error = errors.New("an invalid command was entered")
)

// CreateFile проверяет наличие файла в текущей директории.
// Если его нет, то он будет создан.
func CreateFile() error {
	_, err := os.Stat(tasksPath)
	if err == nil {
		//fmt.Printf("Файл '%s' существует\n", tasksPath)
		return nil

	} else if os.IsNotExist(err) {
		file, err := os.Create(tasksPath)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}

		err = file.Close()
		if err != nil {
			return fmt.Errorf("close file: %w", err)
		}

		fmt.Printf("File '%s' created\n", tasksPath)

	} else {
		return fmt.Errorf("file check: %w", err)
	}

	return nil
}

// GetAllTasks реализует считывание всех задач из файла, преобразует их из JSON в объекты типа Task и возвращает их.
func GetAllTasks() ([]models.Task, error) {
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	allTasks := make([]models.Task, 0)
	if len(data) == 0 {
		return allTasks, nil
	}

	err = json.Unmarshal(data, &allTasks)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return allTasks, nil
}

// printTasks реализует вывод в терминал список задач с преобразованием их статуса в читаемый вид.
func printTasks(tasks []models.Task) error {
	var resBuild strings.Builder
	for _, task := range tasks {
		var status string
		switch task.Status {
		case models.StatusDone:
			status = "Done"
		case models.StatusInProgress:
			status = "In progress"
		case models.StatusNotDone:
			status = "Not started"
		default:
			status = "Incorrect task status"
		}
		resBuild.WriteString(fmt.Sprintf("Index: %d\tName: %s\tStatus: %s\n", task.Index, task.Name, status))
	}

	cmd := exec.Command("less")
	cmd.Stdin = strings.NewReader(resBuild.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// addToFile преобразует полученные объекты типа Task и добавляет обновленный список
// пользовательских задач в созданный файл.
func addToFile(allTasks []models.Task) error {
	tasksJSON, err := json.MarshalIndent(allTasks, "", "\t")
	if err != nil {
		return errors.New("error serializing to JSON")
	}

	err = os.WriteFile(tasksPath, []byte(tasksJSON), 0644)
	if err != nil {
		return errors.New("error writing to file")
	}

	return nil
}

// Add является методом объекта типа Task и реализует добавление новой задачи в список задач. После чего возвращает
// сообщение о результате действия или ошибку.
func Add(tasks *[]models.Task, elements []string) (string, error) {
	newName := strings.ReplaceAll(strings.ReplaceAll(elements[0], "\"", ""), "'", "")
	if newName == "" {
		return "", ErrNameNotExists
	}

	newTask := models.Task{
		Index:  (*tasks)[len(*tasks)-1].Index + 1,
		Name:   newName,
		Status: models.StatusNotDone,
	}

	*tasks = append(*tasks, newTask)

	err := addToFile(*tasks)
	if err != nil {
		return "", fmt.Errorf("addToFile: %w", err)
	}

	return "Task added", nil
}

// Update является методом объекта типа Task и реализует обновление имени и статуса задачи
// по указанному пользователем индексу задачи. После чего возвращает сообщение о результате
// действия или ошибку.
func Update(tasks *[]models.Task, elements []string) (string, error) {
	if elements[0] == "" {
		return "", ErrIndexNotExists
	}

	index, err := strconv.Atoi(elements[0])
	if err != nil {
		return "", errAtoi
	}

	for i, task := range *tasks {
		if task.Index == index {
			if elements[2] == "" {
				return "", ErrStatusNotExists
			}
			status, err := strconv.Atoi(elements[2])
			if err != nil {
				return "", errAtoi
			}
			if status < 0 || status > 2 {
				return "", errIncorrectStatus
			}

			if elements[1] == "" {
				return "", ErrNameNotExists
			}

			(*tasks)[i].Name = elements[1]
			(*tasks)[i].Status = models.TaskStatus(status)

			err = addToFile(*tasks)
			if err != nil {
				return "", fmt.Errorf("addToFile: %w", err)
			}

			return "Task updated", nil
		}
	}

	return "Task not found", nil
}

// Delete является методом объекта типа Task и реализует удаление задачи по индексу из общего списка задач.
// После чего возвращает сообщение о результате действия или ошибку.
func Delete(tasks *[]models.Task, elements []string) (string, error) {
	if elements[0] == "" {
		return "", ErrIndexNotExists
	}

	index, err := strconv.Atoi(elements[0])
	if err != nil {
		return "", errAtoi
	}

	for i, task := range *tasks {
		if task.Index == index {
			if i == len(*tasks)-1 {
				*tasks = (*tasks)[:i]
			} else {
				*tasks = append((*tasks)[:i], (*tasks)[i+1:]...)
			}

			err = addToFile(*tasks)
			if err != nil {
				return "", fmt.Errorf("addToFile: %w", err)
			}

			return "Task deleted", nil
		}
	}

	return "Task not found", nil
}

// UpdateStatus является методом объекта типа Task и реализует обновление статуса задачи
// по указанному пользователем индексу задачи. После чего возвращает сообщение о результате
// действия или ошибку.
func UpdateStatus(tasks *[]models.Task, elements []string) (string, error) {
	if elements[0] == "" {
		return "", ErrIndexNotExists
	}

	index, err := strconv.Atoi(elements[0])
	if err != nil {
		return "", errAtoi
	}

	for i, task := range *tasks {
		if task.Index == index {
			if elements[1] == "" {
				return "", ErrStatusNotExists
			}
			status, err := strconv.Atoi(elements[1])
			if err != nil {
				return "", errAtoi
			}
			if status < 0 || status > 2 {
				return "", errIncorrectStatus
			}

			(*tasks)[i].Status = models.TaskStatus(status)

			err = addToFile(*tasks)
			if err != nil {
				return "", fmt.Errorf("addToFile: %w", err)
			}

			return "Task updated", nil
		}
	}

	return "Task not found", nil
}

// AllTasks передает в функцию для вывода в терминал список всех существующих задач пользователя.
func AllTasks(tasks *[]models.Task) error {
	err := printTasks(*tasks)
	if err != nil {
		return fmt.Errorf("printTasks: %w", err)
	}
	return nil
}

// DoneTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "Выполнено".
func DoneTasks(tasks *[]models.Task) error {
	var result []models.Task

	for _, task := range *tasks {
		if task.Status == models.StatusDone {
			result = append(result, task)
		}
	}

	err := printTasks(result)
	if err != nil {
		return fmt.Errorf("printTasks: %w", err)
	}
	return nil
}

// NotDoneTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "Не начато".
func NotDoneTasks(tasks *[]models.Task) error {
	var result []models.Task

	for _, task := range *tasks {
		if task.Status == models.StatusNotDone {
			result = append(result, task)
		}
	}

	err := printTasks(result)
	if err != nil {
		return fmt.Errorf("printTasks: %w", err)
	}
	return nil
}

// InProgressTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "В процессе".
func InProgressTasks(tasks *[]models.Task) error {
	var result []models.Task

	for _, task := range *tasks {
		if task.Status == models.StatusInProgress {
			result = append(result, task)
		}
	}

	err := printTasks(result)
	if err != nil {
		return fmt.Errorf("printTasks: %w", err)
	}
	return nil
}

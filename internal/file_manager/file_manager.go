package filemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/models"
)

const (
	tasksPath = "tasks.json"
)

var (
	ErrInputElementsCount error = errors.New("передано некорректное количество аргументов")
	errAtoi               error = errors.New("передано некорректное число")
	errIncorrectStatus    error = errors.New("передан некорректный статус задачи")
)

func CreateFile() error {
	_, err := os.Stat(tasksPath)
	if err == nil {
		fmt.Printf("Файл '%s' существует\n", tasksPath)
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

		fmt.Printf("Файл '%s' создан\n", tasksPath)

	} else {
		return fmt.Errorf("проверка файла: %w", err)
	}

	return nil
}

// getAllTasks реализует считывание всех задач из файла, преобразует их из JSON в объекты типа Task и возвращает их.
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
func printTasks(tasks []models.Task) string {
	var result string
	for _, task := range tasks {
		var status string
		switch task.Status {
		case models.StatusDone:
			status = "Выполнено"
		case models.StatusInProgress:
			status = "В процессе"
		case models.StatusNotDone:
			status = "Не начато"
		default:
			status = "Некорректный статус задачи"
		}
		result += fmt.Sprintf("Index: %d\tName: %s\tStatus: %s\n", task.Index, task.Name, status)
	}

	return result
}

// addToFile преобразует полученные объекты типа Task и добавляет обновленный список
// пользовательских задач в созданный файл.
func addToFile(allTasks []models.Task) error {
	tasksJSON, err := json.MarshalIndent(allTasks, "", "\t")
	if err != nil {
		return errors.New("ошибка при сериализации в JSON")
	}

	err = os.WriteFile(tasksPath, []byte(tasksJSON), 0644)
	if err != nil {
		return errors.New("ошибка при записи в файл")
	}

	return nil
}

// Add является методом объекта типа Task и реализует добавление новой задачи в список задач. После чего возвращает
// сообщение о результате действия.
func Add(tasks *[]models.Task, elements []string) (string, error) {
	newName := strings.ReplaceAll(strings.ReplaceAll(elements[0], "\"", ""), "'", "")
	if newName == "" {
		return "", errors.New("введено пустое имя задачи")
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

	return "Задача добавлена", nil
}

// Update является методом объекта типа Task и реализует обновление имени и статуса задачи
// по указанному пользователем индексу задачи. После чего возвращает сообщение о результате действия.
func Update(tasks *[]models.Task, elements []string) (string, error) {
	index, err := strconv.Atoi(elements[0])
	if err != nil {
		return "", errAtoi
	}

	for i, task := range *tasks {
		if task.Index == index {
			status, err := strconv.Atoi(elements[2])
			if err != nil {
				return "", errAtoi
			}
			if status < 0 || status > 2 {
				return "", errIncorrectStatus
			}

			(*tasks)[i].Name = elements[1]
			(*tasks)[i].Status = models.TaskStatus(status)

			err = addToFile(*tasks)
			if err != nil {
				return "", fmt.Errorf("addToFile: %w", err)
			}

			return "Задача обновлена", nil
		}
	}

	return "Задача не найдена", nil
}

// Delete является методом объекта типа Task и реализует удаление задачи по индексу из общего списка задач.
// После чего возвращает сообщение о результате действия.
func Delete(tasks *[]models.Task, elements []string) (string, error) {
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

			return "Задача удалена", nil
		}
	}

	return "Задача не найдена", nil
}

// UpdateStatus является методом объекта типа Task и реализует обновление статуса задачи
// по указанному пользователем индексу задачи. После чего возвращает сообщение о результате действия.
func UpdateStatus(tasks *[]models.Task, elements []string) (string, error) {
	index, err := strconv.Atoi(elements[0])
	if err != nil {
		return "", errAtoi
	}

	for i, task := range *tasks {
		if task.Index == index {
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

			return "Задача обновлена", nil
		}
	}

	return "Задача не найдена", nil
}

// AllTasks передает в функцию для вывода в терминал список всех существующих задач пользователя.
func AllTasks(tasks *[]models.Task) string {
	return printTasks(*tasks)
}

// DoneTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "Выполнено".
func DoneTasks(tasks *[]models.Task) string {
	var result []models.Task

	for _, task := range *tasks {
		if task.Status == models.StatusDone {
			result = append(result, task)
		}
	}
	return printTasks(result)
}

// NotDoneTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "Не начато".
func NotDoneTasks(tasks *[]models.Task) string {
	var result []models.Task

	for _, task := range *tasks {
		if task.Status == models.StatusNotDone {
			result = append(result, task)
		}
	}
	return printTasks(result)
}

// InProgressTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "В процессе".
func InProgressTasks(tasks *[]models.Task) string {
	var result []models.Task

	for _, task := range *tasks {
		if task.Status == models.StatusInProgress {
			result = append(result, task)
		}
	}
	return printTasks(result)
}

// Package taskstorage реализует основное взаимодействие с сущностью "Task".
// В пакете реализован основной функционал вводимых команд пользователем.
// Таким образом, реализованы сериализация и десериализация JSON формата для данных,
// возможность записи в файл и чтение из него.
package taskstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TaskStatus int

// Task структура позволяет воздавать объекты этого типа и преобразовывать их в формат JSON.
type Task struct {
	Index  int        `json:"index"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
}

const (
	StatusNotDone TaskStatus = iota
	StatusInProgress
	StatusDone
	TasksPath = "tasks.json"
)

var errInputElementsCount error = errors.New("передано некорректное количество аргументов")
var errAtoi error = errors.New("введено некорректное число")

// printTasks реализует вывод в терминал список задач с преобразованием их статуса в читаемый вид.
func printTasks(tasks []Task) string {
	var result string
	for _, task := range tasks {
		var status string
		switch task.Status {
		case StatusDone:
			status = "Выполнено"
		case StatusInProgress:
			status = "В процессе"
		case StatusNotDone:
			status = "Не начато"
		default:
			status = "Некорректный статус задачи"
		}
		result += fmt.Sprintf("Index: %d\nName: %s\nStatus: %s\n", task.Index, task.Name, status)
	}
	return result
}

// addToFile преобразует полученные объекты типа Task и добавляет обновленный список
// пользовательских задач в созданный файл.
func addToFile(allTasks []Task) error {
	tasksJSON, err := json.MarshalIndent(allTasks, "", "\t")
	if err != nil {
		return errors.New("ошибка при сериализации в JSON")
	}

	err = os.WriteFile(TasksPath, []byte(tasksJSON), 0644)
	if err != nil {
		return errors.New("ошибка при записи в файл")
	}

	return nil
}

// getAllTasks реализует считывание всех задач из файла, преобразует их из JSON в объекты типа Task и возвращает их.
func getAllTasks() ([]Task, error) {
	data, err := os.ReadFile(TasksPath)
	if err != nil {
		return nil, errors.New("ошибка при чтении файла")
	}

	var allTasks []Task
	if len(data) == 0 {
		return allTasks, nil
	}

	err = json.Unmarshal(data, &allTasks)
	if err != nil {
		return nil, errors.New("ошибка при десериализации из JSON")
	}

	return allTasks, nil
}

// Add является методом объекта типа Task и реализует добавление новой задачи в список задач. После чего возвращает
// сообщение о результате действия.
func (t Task) Add(elements []string) (string, error) {
	if len(elements) != 1 {
		return "", errInputElementsCount
	}

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	t.Index = len(allTasks) + 1
	elements[0] = strings.ReplaceAll(elements[0], "\"", "")
	t.Name = strings.ReplaceAll(elements[0], "'", "")
	if t.Name == "" {
		return "", errors.New("введено пустое имя задачи")
	}

	allTasks = append(allTasks, t)

	err = addToFile(allTasks)
	if err != nil {
		return "", err
	}

	return "Задача добавлена", nil
}

// Update является методом объекта типа Task и реализует обновление имени и статуса задачи
// по указанному пользователем индексу задачи. После чего возвращает сообщение о результате действия.
func (t Task) Update(elements []string) (string, error) {
	if len(elements) != 3 {
		return "", errInputElementsCount
	}

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(allTasks); i++ {
		index, err := strconv.Atoi(elements[0])
		if err != nil {
			return "", errAtoi
		}

		if allTasks[i].Index == index {
			status, err := strconv.Atoi(elements[2])
			if err != nil {
				return "", errAtoi
			}

			t.Index = index
			t.Name, allTasks[i].Name = elements[1], elements[1]
			t.Status, allTasks[i].Status = TaskStatus(status), TaskStatus(status)

			err = addToFile(allTasks)
			if err != nil {
				return "", err
			}

			return printTasks([]Task{t}), nil
		}
	}

	return "Задача не найдена", nil
}

// Delete является методом объекта типа Task и реализует удаление задачи по индексу из общего списка задач.
// После чего возвращает сообщение о результате действия.
func (t Task) Delete(elements []string) (string, error) {
	if len(elements) != 1 {
		return "", errInputElementsCount
	}

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(allTasks); i++ {
		index, err := strconv.Atoi(elements[0])
		if err != nil {
			return "", errAtoi
		}

		if allTasks[i].Index == index {
			allTasks = append(allTasks[:i], allTasks[i+1:]...)
			err = addToFile(allTasks)
			if err != nil {
				return "", err
			}

			return printTasks(allTasks), nil
		}
	}

	return "Задача не найдена", nil
}

// UpdateStatus является методом объекта типа Task и реализует обновление статуса задачи
// по указанному пользователем индексу задачи. После чего возвращает сообщение о результате действия.
func (t Task) UpdateStatus(elements []string) (string, error) {
	if len(elements) != 2 {
		return "", errInputElementsCount
	}

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(allTasks); i++ {
		index, err := strconv.Atoi(elements[0])
		if err != nil {
			return "", errAtoi
		}

		if allTasks[i].Index == index {
			status, err := strconv.Atoi(elements[1])
			if err != nil {
				return "", errAtoi
			}

			t.Index = index
			t.Name = allTasks[i].Name
			t.Status, allTasks[i].Status = TaskStatus(status), TaskStatus(status)

			err = addToFile(allTasks)
			if err != nil {
				return "", err
			}

			return printTasks([]Task{t}), nil
		}
	}
	return "Задача не найдена", nil
}

// AllTasks передает в функцию для вывода в терминал список всех существующих задач пользователя.
func AllTasks() (string, error) {
	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}
	return printTasks(allTasks), nil
}

// DoneTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "Выполнено".
func DoneTasks() (string, error) {
	var result []Task

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	for _, task := range allTasks {
		if task.Status == StatusDone {
			result = append(result, task)
		}
	}
	return printTasks(result), nil
}

// NotDoneTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "Не начато".
func NotDoneTasks() (string, error) {
	var result []Task

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	for _, task := range allTasks {
		if task.Status == StatusNotDone {
			result = append(result, task)
		}
	}
	return printTasks(result), nil
}

// InProgressTasks передает в функцию для вывода в терминал список всех существующих задач пользователя
// со статусом "В процессе".
func InProgressTasks() (string, error) {
	var result []Task

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	for _, task := range allTasks {
		if task.Status == StatusInProgress {
			result = append(result, task)
		}
	}
	return printTasks(result), nil
}

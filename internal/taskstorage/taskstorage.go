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
		}
		result += fmt.Sprintf("Index: %d\nName: %s\nStatus: %s\n", task.Index, task.Name, status)
	}
	return result
}

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

func (t Task) Add(elements []string) (string, error) {
	if len(elements) != 1 {
		return "", errInputElementsCount
	}

	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}

	t.Index = len(allTasks) + 1
	elements[0] = strings.ReplaceAll(elements[0],"\"", "")
	t.Name = strings.ReplaceAll(elements[0],"'", "")
	if t.Name == ""{
		return "", errors.New("введено пустое имя задачи")
	}

	allTasks = append(allTasks, t)

	err = addToFile(allTasks)
	if err != nil {
		return "", err
	}

	return "Задача добавлена", nil
}

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

func AllTasks() (string, error) {
	allTasks, err := getAllTasks()
	if err != nil {
		return "", err
	}
	return printTasks(allTasks), nil
}

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

package taskstorage

import (
	"encoding/json"
	"errors"
	"fmt"

	// "fmt"
	"log"
	"os"
	"strconv"
)

type TaskStatus int

const (
	TasksPath   = "tasks.json"
	CounterPath = "count.txt"
)

const (
	StatusNotDone TaskStatus = iota
	StatusInProgress
	StatusDone
)

type Task struct {
	Index  int        `json:"index"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
}

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

		result += fmt.Sprintf(`Index: %d
Name: %s
Status: %s

`, task.Index, task.Name, status)
	}
	return result
}

func (t Task) Add(elements []string) (string, error) {
	if len(elements) != 1 {
		return "", errors.New("передано некорректное количество аргументов")
	}

	countStr, err := os.ReadFile(CounterPath)
	if err != nil {
		return "", err
	}
	count, err := strconv.Atoi(string(countStr))
	if err != nil {
		return "", err
	}
	t.Index = count + 1
	t.Name = elements[0]

	file, err := os.OpenFile(TasksPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	if stat.Size() == 0 {
		// Файл пустой - пишем начальный массив
		_, err = file.WriteString("[\n\t")
		if err != nil {
			return "", err
		}
	} else {
		// Перемещаемся перед закрывающей скобкой
		file.Seek(-1, 2) // Последний символ
		lastChar := make([]byte, 1)
		file.Read(lastChar)

		if lastChar[0] == ']' {
			// Массив не пустой - добавляем запятую и новую строку
			file.Seek(-1, 2)
			_, err = file.WriteString(",\n\t")
			if err != nil {
				return "", err
			}
		}
	}

	taskJSON, err := json.MarshalIndent(t, "\t", "\t")
	if err != nil {
		return "", err
	}

	_, err = file.Write(taskJSON)
	if err != nil {
		return "", err
	}

	_, err = file.WriteString("\n]")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(CounterPath, []byte(strconv.Itoa(t.Index)), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return "Задача добавлена", err
}

func (t Task) Update(elements []string) {
	//TODO
}

func (t Task) Delete(elements []string) {
	//TODO
}

func (t Task) UpdateStatus(elements []string) {
	//TODO
}

func (t Task) AllTasks() (string, error) {
	var tasks []Task

	data, err := os.ReadFile(TasksPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &tasks)
	return printTasks(tasks), nil
}

func (t Task) DoneTasks() (string, error) {
	var tasks []Task
	var result []Task

	data, err := os.ReadFile(TasksPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &tasks)
	for _, task := range tasks {
		if task.Status == StatusDone {
			result = append(result, task)
		}
	}
	return printTasks(result), nil
}

func (t Task) NotDoneTasks() (string, error) {
	var tasks []Task
	var result []Task

	data, err := os.ReadFile(TasksPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &tasks)
	for _, task := range tasks {
		if task.Status == StatusNotDone {
			result = append(result, task)
		}
	}
	return printTasks(result), nil
}

func (t Task) InProgressTasks() (string, error) {
	var tasks []Task
	var result []Task

	data, err := os.ReadFile(TasksPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &tasks)
	for _, task := range tasks {
		if task.Status == StatusInProgress {
			result = append(result, task)
		}
	}
	return printTasks(result), nil
}

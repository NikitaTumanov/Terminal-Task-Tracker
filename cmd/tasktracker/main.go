package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/taskstorage"
)

func Read() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\r\n")
	return input
}

func splitInput(input string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false

	for _, char := range input {
		switch {
		case char == '"':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case char == ' ' && !inQuotes:
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func main() {
	fmt.Println("Task Manager Started")

	_, err := os.Stat(taskstorage.TasksPath)
	if err == nil {
		fmt.Printf("Файл '%s' существует\n", taskstorage.TasksPath)
	} else if os.IsNotExist(err) {
		file, err := os.Create(taskstorage.TasksPath)
		if err != nil {
			panic(err)
		}
		file.Close()
		fmt.Printf("Файл '%s' создан\n", taskstorage.TasksPath)
	} else {
		fmt.Printf("Ошибка при проверке файла: %v\n", err)
		panic(err)
	}

	_, err = os.Stat(taskstorage.CounterPath)
	if err == nil {
		fmt.Printf("Файл '%s' существует\n", taskstorage.CounterPath)
	} else if os.IsNotExist(err) {
		file, err := os.Create(taskstorage.CounterPath)
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString("0")
		if err != nil {
			panic(err)
		}
		file.Close()
		fmt.Printf("Файл '%s' создан\n", taskstorage.CounterPath)
	} else {
		fmt.Printf("Ошибка при проверке файла: %v\n", err)
		panic(err)
	}

	var task taskstorage.Task
	for {
		fmt.Print("Enter text: ")
		input := Read()

		elements := splitInput(input)

		switch strings.ToLower(elements[0]) {
		case "add":
			result, err := task.Add(elements[1:])
			if err != nil {
				log.Fatal(err)
				fmt.Println("Что-то не так")
				continue
			}
			fmt.Println(result)
		case "update":
			fmt.Println("Функционал еще не реализован")
		case "delete":
			fmt.Println("Функционал еще не реализован")
		case "updatestatus":
			fmt.Println("Функционал еще не реализован")
		case "alltasks":
			t, _ := task.AllTasks()
			fmt.Println(t)
		case "donetasks":
			t, _ := task.DoneTasks()
			fmt.Println(t)
		case "notdonetasks":
			t, _ := task.NotDoneTasks()
			fmt.Println(t)
		case "inprogresstasks":
			t, _ := task.InProgressTasks()
			fmt.Println(t)
		case "help":
			fmt.Println(`Add <Task name>'>
Update <Task Index> <Task Status>
	Task Statuses:
		0 - 
Delete
UpdateStatus
AllTasks
DoneTasks
NotDoneTasks
InProgressTasks
Help
Exit`)
		case "exit":
			return
		default:
			fmt.Println("Введена некорректная команда")
		}
	}
}

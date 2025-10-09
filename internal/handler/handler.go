// Package handler реализует обработку команд, которые поступают от пользователя.
// Команда разбивается на атрибуты и передается в соответствующую функцию в зависимости от ожидаемого результата.
package handler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/taskstorage"
)

// Read выполняет чтение команд из терминала до тех пор, пока ввод будет пустым,
// и возвращает прочитанную команду в виде строки.
// При вводе пустой строки в терминал выведется подсказка для пользователя.
func Read() string {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimRight(input, "\r\n")

		if input != "" {
			return input
		}

		fmt.Println("Чтобы получить информацию о доступных командах, введите help")
	}
}

// splitInput получает на вход считанную команду пользователя в виде строки
// и разбивает ее на составляющие (команда и аргументы), после чего возвращает их в виде слайса типа строки.
// Ключевой особенностью является то, что символы, заключенные в кавычки считаются одним аргументом.
func splitInput(input string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, char := range input {
		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = rune(0)
			}
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

// Handle вызывает функцию считывания пользовательского ввода до тех пор, пока не поступит команда Exit.
// В иных случаях функция вызывает соответствующий метод в зависимости от команды пользователя
// и выводит результат в терминал.
func Handle() {
	var task taskstorage.Task

	for {
		fmt.Print("Введите команду: ")
		input := Read()
		elements := splitInput(input)

		switch strings.ToLower(elements[0]) {
		case "add":
			result, err := task.Add(elements[1:])
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "update":
			result, err := task.Update(elements[1:])
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "delete":
			result, err := task.Delete(elements[1:])
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "updatestatus":
			result, err := task.UpdateStatus(elements[1:])
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "alltasks":
			result, err := taskstorage.AllTasks()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "donetasks":
			result, err := taskstorage.DoneTasks()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "notdonetasks":
			result, err := taskstorage.NotDoneTasks()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "inprogresstasks":
			result, err := taskstorage.InProgressTasks()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(result)
		case "help":
			fmt.Println(`	Add "<Task name>"
	Update <Task Index> "<New Task Name>" <New Task Status>
		Task Statuses:
			0 - Не начато
			1 - В процессе
			2 - Выполнено
	Delete <Task Index>
	UpdateStatus <Task Index> <New Task Status>
		Task Statuses:
			0 - Не начато
			1 - В процессе
			2 - Выполнено
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

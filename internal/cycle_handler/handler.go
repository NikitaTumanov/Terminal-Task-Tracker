// Package handler реализует обработку команд, которые поступают от пользователя.
// Команда разбивается на атрибуты и передается в соответствующую функцию в зависимости от ожидаемого результата.
package cyclehandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

// Read выполняет чтение команд из терминала до тех пор, пока ввод будет пустым,
// и возвращает прочитанную команду в виде строки.
// При вводе пустой строки в терминал выведется подсказка для пользователя.
func read() string {
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
func (s *storage) Handle() error {
	fmt.Println("Task Manager Started")
	for {
		fmt.Print("Введите команду: ")
		input := read()
		elements := splitInput(input)

		switch strings.ToLower(elements[0]) {
		case "add":
			if len(elements) != 2 {
				fmt.Println(filemanager.ErrInputElementsCount)
				continue
			}
			result, err := filemanager.Add(&s.tasks, elements[1:])
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(result)

		case "update":
			if len(elements) != 4 {
				fmt.Println(filemanager.ErrInputElementsCount)
				continue
			}
			result, err := filemanager.Update(&s.tasks, elements[1:])
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(result)

		case "delete":
			if len(elements) != 2 {
				fmt.Println(filemanager.ErrInputElementsCount)
				continue
			}
			result, err := filemanager.Delete(&s.tasks, elements[1:])
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(result)

		case "updatestatus":
			if len(elements) != 3 {
				fmt.Println(filemanager.ErrInputElementsCount)
				continue
			}
			result, err := filemanager.UpdateStatus(&s.tasks, elements[1:])
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(result)

		case "alltasks":
			err := filemanager.AllTasks(&s.tasks)
			if err != nil {
				fmt.Println(err)
			}

		case "donetasks":
			err := filemanager.DoneTasks(&s.tasks)
			if err != nil {
				fmt.Println(err)
			}

		case "notdonetasks":
			err := filemanager.NotDoneTasks(&s.tasks)
			if err != nil {
				fmt.Println(err)
			}

		case "inprogresstasks":
			err := filemanager.InProgressTasks(&s.tasks)
			if err != nil {
				fmt.Println(err)
			}

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
			return nil
		default:
			fmt.Println("Введена некорректная команда")
		}
	}
}

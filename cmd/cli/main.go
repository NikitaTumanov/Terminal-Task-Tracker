// Package main реализует точку запуска приложения.
// В этом пакете создается .json файл (если его нет), в котором будут храниться созданные задачи пользователя.
// После чего вызывается обработчик пользовательских команд.
package main

import (
	"flag"
	"fmt"
	"strings"

	cyclehandler "github.com/NikitaTumanov/terminalTaskTracker/internal/cycle_handler"
	filemanager "github.com/NikitaTumanov/terminalTaskTracker/internal/file_manager"
	flaghandler "github.com/NikitaTumanov/terminalTaskTracker/internal/flag_handler"
)

type handler interface {
	Handle() error
	Update()
}

// main запускает работу приложения.
// Функция создает в той же директории JSON файл для записи задач и вызывает метод обработки команд.
// При возникновении ошибки при работе с файлом приложение прекращает работу.
func main() {
	var (
		handler    handler
		mode       = flag.String("mode", "flag", "mode")
		command    string
		taskIndex  = flag.Int("index", -1, "index")
		taskName   = flag.String("name", "", "name")
		taskStatus = flag.String("status", "", "status")
	)
	flag.StringVar(&command, "c", "", "command")
	flag.Parse()

	err := filemanager.CreateFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	if command == "" {
		handler, err = cyclehandler.New()
		if err != nil {
			fmt.Println(err)
		}

	} else if strings.EqualFold(strings.ToLower(*mode), "flag") {
		handler, err = flaghandler.New(*taskIndex, command, *taskName, *taskStatus)
		if err != nil {
			fmt.Println(err)
		}

	}

	handler.Update()
	err = handler.Handle()
	if err != nil {
		fmt.Println(err)
	}
}

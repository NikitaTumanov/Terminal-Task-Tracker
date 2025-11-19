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
}

// main запускает работу приложения.
// Функция создает в той же директории JSON файл для записи задач и вызывает метод обработки команд.
// При возникновении ошибки при работе с файлом приложение прекращает работу.
func main() {
	fmt.Println("Task Manager Started")
	var handler handler

	err := filemanager.CreateFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	var mode = flag.String("mode", "flag", "mode")
	flag.Parse()

	if strings.EqualFold(strings.ToLower(*mode), "interactive") {
		handler, err = cyclehandler.New()
		if err != nil {
			fmt.Println(err)
		}

	} else if strings.EqualFold(strings.ToLower(*mode), "flag") {
		handler, err = flaghandler.New()
		if err != nil {
			fmt.Println(err)
		}

	} else {
		err := fmt.Errorf("введено некорректное значение режима работы mode: %s", *mode)
		fmt.Println(err, "\nОжидается: --mode: interactive|flag(по умолчанию)")
		return
	}

	err = handler.Handle()
	if err != nil {
		fmt.Println(err)
	}
}

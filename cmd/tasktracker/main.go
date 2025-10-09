// Package main реализует точку запуска приложения.
// В этом пакете создается .json файл (если его нет), в котором будут храниться созданные задачи пользователя.
// После чего вызывается обработчик пользовательских команд.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/handler"
	"github.com/NikitaTumanov/terminalTaskTracker/internal/taskstorage"
)

// main запускает работу приложения.
// Функция создает в той же директории JSON файл для записи задач и вызывает метод обработки команд.
// При возникновении ошибки при работе с файлом приложение прекращает работу.
func main() {
	fmt.Println("Task Manager Started")

	_, err := os.Stat(taskstorage.TasksPath)
	if err == nil {
		log.Printf("Файл '%s' существует\n", taskstorage.TasksPath)
	} else if os.IsNotExist(err) {
		file, err := os.Create(taskstorage.TasksPath)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}
		log.Printf("Файл '%s' создан\n", taskstorage.TasksPath)
	} else {
		log.Fatalf("Ошибка при проверке файла: %v\n", err)
	}

	handler.Handle()
}

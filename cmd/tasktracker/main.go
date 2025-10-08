package main

import (
	"fmt"
	"os"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/handler"
	"github.com/NikitaTumanov/terminalTaskTracker/internal/taskstorage"
)

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

	handler.Handle()
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/handler"
	"github.com/NikitaTumanov/terminalTaskTracker/internal/taskstorage"
)

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
		file.Close()
		log.Printf("Файл '%s' создан\n", taskstorage.TasksPath)
	} else {
		log.Fatalf("Ошибка при проверке файла: %v\n", err)
	}

	handler.Handle()
}

package main

import (
	"fmt"
	"strings"

	"github.com/NikitaTumanov/terminalTaskTracker/internal/inputreader"
)

func main() {
	fmt.Println("Task Manager Started")

	
	reader := inputreader.Reader{}
	for strings.ToLower(reader.Input) != "exit\n" {
		fmt.Println("Enter text: ")
		reader.Read()
		fmt.Println(reader.Input)
	}

}

package models

import "time"

type TaskStatus int

// Task структура описывает сущность Task.
type Task struct {
	Index  int        `json:"index"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
}

const (
	StatusNotDone TaskStatus = iota
	StatusInProgress
	StatusDone
	TimeOut = time.Second * 5
)

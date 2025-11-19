package models

type TaskStatus int

// Task структура позволяет воздавать объекты этого типа и преобразовывать их в формат JSON.
type Task struct {
	Index  int        `json:"index"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
}

const (
	StatusNotDone TaskStatus = iota
	StatusInProgress
	StatusDone
)

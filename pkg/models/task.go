package models

import "fmt"

type Task struct {
	Id          int64  `json:"id,omitempty"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	IsCompleted bool   `json:"is_completed"`
}

func (t *Task) GetInfo() string {
	return fmt.Sprintf("Id: %d Title: %s Text: %s Completed: %t", t.Id, t.Title, t.Text, t.IsCompleted)
}

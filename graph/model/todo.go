package model

import (
	"fmt"

	"github.com/maruware/gqlgen-todos/entity"
)

type Todo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID string `json:"user"`
}

func NewTodoFromEntity(e *entity.Todo) *Todo {
	return &Todo{
		ID:     fmt.Sprintf("%d", e.ID),
		Text:   e.Text,
		Done:   e.Done,
		UserID: fmt.Sprintf("%d", e.UserID),
	}
}

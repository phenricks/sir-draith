package entities

import "time"

// Task representa uma tarefa no sistema
type Task struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Status      string    `bson:"status"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// NewTask cria uma nova tarefa com valores padrão
func NewTask(title, description string) *Task {
	now := time.Now()
	return &Task{
		Title:       title,
		Description: description,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

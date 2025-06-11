package tasks

import "time"

const (
	StatusCreated   = "created"   // Задача создана
	StatusRunning   = "running"   // Задача выполняется
	StatusCompleted = "completed" // Задача завершена
	StatusFailed    = "failed"    // Ошибка при выполнении задачи
)

type Task struct {
	ID          uint64        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	FinishedAt  time.Time     `json:"finished_at"`
	Duration    time.Duration `json:"duration"`
	Status      string        `json:"status"`
}

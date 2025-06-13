package payloads

import "time"

const (
	StatusCreated   = "created"   // Задача создана
	StatusRunning   = "running"   // Задача выполняется
	StatusCompleted = "completed" // Задача завершена
	StatusFailed    = "failed"    // Ошибка при выполнении задачи
)

type Task struct {
	ID          uint64        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	FinishedAt  time.Time     `json:"finished_at"`
	Duration    time.Duration `json:"duration"`
	Status      string        `json:"status"`
	Error       string        `json:"error"`
	DeleteChan  chan struct{} `json:"-"`
}

func (task *Task) SetDuration() {
	var duration float64
	if task.Status == StatusCompleted || task.Status == StatusFailed {
		duration = task.FinishedAt.Sub(task.CreatedAt).Seconds()
	} else {
		duration = time.Now().Sub(task.CreatedAt).Seconds()
	}
	task.Duration = time.Duration(duration)
}

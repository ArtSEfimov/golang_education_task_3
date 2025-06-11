package service

import (
	"io_bound_task/internal/tasks"
	"time"
)

const (
	StatusIdle    = iota // Процессор в простое
	StatusWorking        // Процессор работает
)

type machine struct {
	boundProcessor *Processor
	inputTaskChan  chan *tasks.Task
	outputTaskChan chan *tasks.Task
	status         int
}

func newMachine(processor *Processor) *machine {
	return &machine{
		boundProcessor: processor,
		status:         StatusIdle,
	}
}

func (m *machine) assignTask(task *tasks.Task) {
	m.status = StatusWorking
	m.boundProcessor.moveToBusyWorkerPool(m)
	m.inputTaskChan <- task
	m.processTask()
}

func (m *machine) processTask() {
	task := <-m.inputTaskChan
	task.Status = tasks.StatusRunning
	time.Sleep(5 * time.Minute)
	m.outputTaskChan <- task
	m.catchCompletedTask()
}

func (m *machine) catchCompletedTask() {
	task := <-m.outputTaskChan
	m.status = StatusIdle
	m.boundProcessor.moveToFreeWorkerPool(m)
	m.boundProcessor.mtx.Lock()
	defer m.boundProcessor.mtx.Unlock()
	task.Status = tasks.StatusCompleted
	task.FinishedAt = time.Now()
	task.Duration = task.FinishedAt.Sub(task.CreatedAt)
}

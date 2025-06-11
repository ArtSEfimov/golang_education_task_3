package service

import (
	"io_bound_task/internal/tasks/payloads"
	"time"
)

const (
	StatusIdle    = iota // Процессор в простое
	StatusWorking        // Процессор работает
)

type machine struct {
	boundProcessor *Processor
	inputTaskChan  chan *payloads.Task
	outputTaskChan chan *payloads.Task
	status         int
}

func newMachine(processor *Processor) *machine {
	return &machine{
		boundProcessor: processor,
		inputTaskChan:  make(chan *payloads.Task),
		outputTaskChan: make(chan *payloads.Task),
		status:         StatusIdle,
	}
}

func (m *machine) assignTask(task *payloads.Task) {
	m.status = StatusWorking
	m.boundProcessor.moveToBusyWorkerPool(m)
	go m.processTask()
	m.inputTaskChan <- task
}

func (m *machine) processTask() {
	task := <-m.inputTaskChan
	task.Status = payloads.StatusRunning
	time.Sleep(5 * time.Minute)
	go m.catchCompletedTask()
	m.outputTaskChan <- task
}

func (m *machine) catchCompletedTask() {
	task := <-m.outputTaskChan
	m.status = StatusIdle
	m.boundProcessor.moveToFreeWorkerPool(m)
	m.boundProcessor.mtx.Lock()
	defer m.boundProcessor.mtx.Unlock()
	task.Status = payloads.StatusCompleted
	task.FinishedAt = time.Now()
	task.SetDuration()
}

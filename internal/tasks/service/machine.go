package service

import (
	"fmt"
	"io_bound_task/internal/tasks/payloads"
	rand2 "math/rand"
	"time"
)

const (
	StatusIdle    = iota // Процессор в простое
	StatusWorking        // Процессор работает
)

const timeSleep = 2

type machine struct {
	boundProcessor *Processor
	inputTaskChan  chan *payloads.Task
	outputTaskChan chan *payloads.Task
	errorTaskChan  chan error
	status         int
}

func newMachine(processor *Processor) *machine {
	return &machine{
		boundProcessor: processor,
		inputTaskChan:  make(chan *payloads.Task),
		outputTaskChan: make(chan *payloads.Task),
		errorTaskChan:  make(chan error),
		status:         StatusIdle,
	}
}

func (m *machine) assignTask(task *payloads.Task) {
	task.Status = payloads.StatusRunning
	m.status = StatusWorking
	m.boundProcessor.moveToBusyWorkerPool(m)
	go m.processTask()
	m.inputTaskChan <- task
}

func (m *machine) processTask() {
	task := <-m.inputTaskChan
	go m.catchCompletedTask()

	select {
	case <-task.DeleteChan:
		m.errorTaskChan <- nil
		m.outputTaskChan <- task
	case <-time.After(time.Duration(timeSleep) * time.Minute):
		// This block simulates a potential error occurring during task processing with a 25% probability
		rnd := rand2.New(rand2.NewSource(time.Now().UnixNano()))
		if rnd.Uint64()%4 == 0 {
			taskError := fmt.Errorf("unable to complete the task")
			m.errorTaskChan <- taskError
		} else {
			m.errorTaskChan <- nil
		}
		m.outputTaskChan <- task
	}
}

func (m *machine) catchCompletedTask() {
	err := <-m.errorTaskChan
	task := <-m.outputTaskChan
	m.status = StatusIdle
	m.boundProcessor.moveToFreeWorkerPool(m)
	m.boundProcessor.mtx.Lock()
	defer m.boundProcessor.mtx.Unlock()

	if err != nil {
		task.Error = err.Error()
		task.Status = payloads.StatusFailed
	} else {
		task.Status = payloads.StatusCompleted
	}
	task.FinishedAt = time.Now()
	task.SetDuration()
}

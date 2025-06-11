package service

import (
	"io_bound_task/internal/tasks"
	"sync"
)

const totalMachines = 10

type workers = map[*machine]struct{}

type Processor struct {
	taskQueue    []*tasks.Task
	freeWorkers  workers
	busyWorkers  workers
	pendingTasks chan *tasks.Task
	mtx          *sync.RWMutex
}

func NewProcessor() *Processor {
	processor := Processor{
		taskQueue:   make([]*tasks.Task, 10),
		freeWorkers: make(workers, 10),
		busyWorkers: make(workers, 10),
	}
	processor.freeWorkers[newMachine(&processor)] = struct{}{}
	return &processor
}

func (processor *Processor) AddTask(task *tasks.Task) {
	processor.taskQueue = append(processor.taskQueue, task)
}

func (processor *Processor) Start() {
	go func() {

		for task := range processor.pendingTasks {
			for {
				freeWorker := processor.getNextFreeWorker()
				if freeWorker != nil {
					freeWorker.assignTask(task)
					break
				}
			}
		}

	}()

	for {
		if len(processor.taskQueue) == 0 {
			continue
		}
		if processor.getFreeWorkersCount() == 0 && processor.getBusyWorkersCount() < totalMachines {
			processor.createFreeWorker()
		}
		for _, task := range processor.taskQueue {
			freeWorker := processor.getNextFreeWorker()
			if freeWorker != nil {
				freeWorker.assignTask(task)
				continue
			}
			processor.pendingTasks <- task
		}
	}
}

func (processor *Processor) createFreeWorker() {
	processor.mtx.Lock()
	defer processor.mtx.Unlock()
	processor.freeWorkers[newMachine(processor)] = struct{}{}
}

func (processor *Processor) getNextFreeWorker() *machine {
	if processor.getFreeWorkersCount() == 0 {
		return nil
	}
	processor.mtx.RLock()
	defer processor.mtx.RUnlock()
	for worker := range processor.freeWorkers {
		return worker
	}
	return nil
}

func (processor *Processor) getFreeWorkersCount() int {
	processor.mtx.RLock()
	defer processor.mtx.RUnlock()
	return len(processor.freeWorkers)
}

func (processor *Processor) getBusyWorkersCount() int {
	processor.mtx.RLock()
	defer processor.mtx.RUnlock()
	return len(processor.busyWorkers)
}

func (processor *Processor) moveToBusyWorkerPool(worker *machine) {
	processor.mtx.Lock()
	defer processor.mtx.Unlock()
	delete(processor.freeWorkers, worker)
	processor.busyWorkers[worker] = struct{}{}
}

func (processor *Processor) moveToFreeWorkerPool(worker *machine) {
	processor.mtx.Lock()
	defer processor.mtx.Unlock()
	delete(processor.busyWorkers, worker)
	processor.freeWorkers[worker] = struct{}{}
}

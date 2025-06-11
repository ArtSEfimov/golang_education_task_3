package service

import (
	"io_bound_task/internal/tasks/payloads"
	"sync"
)

const totalMachines = 10

const (
	initialWorkerPoolCapacity = 10
	initialTaskQueueCapacity  = 10
)

type workers = map[*machine]struct{}

type Processor struct {
	taskQueue       []*payloads.Task
	freeWorkers     workers
	busyWorkers     workers
	pendingTasks    chan *payloads.Task
	mtx             *sync.RWMutex
	taskQueueCond   *sync.Cond
	freeWorkersCond *sync.Cond
}

func NewProcessor() *Processor {
	processor := Processor{
		taskQueue:    make([]*payloads.Task, 0, initialTaskQueueCapacity),
		freeWorkers:  make(workers, initialWorkerPoolCapacity),
		busyWorkers:  make(workers, initialWorkerPoolCapacity),
		pendingTasks: make(chan *payloads.Task),
	}
	mtx := &sync.RWMutex{}
	processor.mtx = mtx
	processor.taskQueueCond = sync.NewCond(mtx)
	processor.freeWorkersCond = sync.NewCond(mtx)
	processor.freeWorkers[newMachine(&processor)] = struct{}{}
	return &processor
}

func (processor *Processor) AddTask(task *payloads.Task) {
	processor.mtx.Lock()
	defer processor.mtx.Unlock()
	processor.taskQueue = append(processor.taskQueue, task)
	processor.taskQueueCond.Signal()
}

func (processor *Processor) Start() {
	go func() {
		for task := range processor.pendingTasks {
			for {
				freeWorker := processor.getNextFreeWorker()
				if freeWorker == nil {
					processor.mtx.Lock()
					processor.freeWorkersCond.Wait()
					processor.mtx.Unlock()
					continue
				}
				freeWorker.assignTask(task)
				break
			}
		}
	}()

	for {
		if processor.getTaskQueueLength() == 0 {
			processor.mtx.Lock()
			processor.taskQueueCond.Wait()
			processor.mtx.Unlock()
			continue
		}
		if processor.getFreeWorkersCount() == 0 && processor.getBusyWorkersCount() < totalMachines {
			processor.createFreeWorker()
		}

		nextTask := processor.getNextTask()
		freeWorker := processor.getNextFreeWorker()
		if freeWorker != nil {
			freeWorker.assignTask(nextTask)
		} else {
			processor.pendingTasks <- nextTask
		}
	}
}

func (processor *Processor) getTaskQueueLength() int {
	processor.mtx.RLock()
	defer processor.mtx.RUnlock()
	return len(processor.taskQueue)
}

func (processor *Processor) getNextTask() *payloads.Task {
	processor.mtx.Lock()
	defer processor.mtx.Unlock()
	task := processor.taskQueue[0]
	processor.taskQueue = processor.taskQueue[1:]
	return task
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
	processor.freeWorkersCond.Signal()
}

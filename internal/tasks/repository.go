package tasks

import (
	"fmt"
	"io_bound_task/internal/tasks/payloads"
	"io_bound_task/internal/tasks/service"
	"slices"
	"time"
)

const initialDBSize = 10

type MapDB = map[uint64]*payloads.Task
type OrderedID = []uint64

type Repository struct {
	idCounter uint64
	DB        MapDB
	Order     OrderedID
}

func NewRepository() *Repository {
	db := make(MapDB, initialDBSize)
	order := make(OrderedID, initialDBSize)
	idCounter := uint64(0)

	return &Repository{
		idCounter: idCounter,
		DB:        db,
		Order:     order,
	}
}

func (repository *Repository) increaseIDCounter() {
	repository.idCounter++
	repository.Order = append(repository.Order, repository.idCounter)
}

func (repository *Repository) getNextID() uint64 {
	repository.increaseIDCounter()
	return repository.idCounter
}

func (repository *Repository) FindByID(id uint64) (*payloads.Task, error) {
	task, ok := repository.DB[id]
	if !ok {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}
	task.SetDuration()
	return task, nil
}

func (repository *Repository) GetAllTasks() (*payloads.AllTasksResponse, error) {
	var allTasksResponse payloads.AllTasksResponse
	for _, task := range repository.DB {
		task.SetDuration()
		allTasksResponse.Tasks = append(allTasksResponse.Tasks, *task)
	}
	return &allTasksResponse, nil
}

func (repository *Repository) GetAllTasksInOrder() (*payloads.AllTasksResponse, error) {
	var allTasksResponse payloads.AllTasksResponse
	for _, id := range repository.Order {
		if task, ok := repository.DB[id]; ok {
			task.SetDuration()
			allTasksResponse.Tasks = append(allTasksResponse.Tasks, *task)
		}
	}
	return &allTasksResponse, nil
}

func (repository *Repository) Delete(id uint64, processor *service.Processor) error {
	task, ok := repository.DB[id]
	if !ok {
		return fmt.Errorf("task with ID %d not found", id)
	}
	if task.Status == "completed" {
	}
	if task.Status == payloads.StatusRunning {
		task.DeleteChan <- struct{}{}
	} else if task.Status == payloads.StatusCreated {
		processor.RemoveTask(task)
	}
	delete(repository.DB, id)
	repository.Order = slices.DeleteFunc(repository.Order, func(v uint64) bool {
		return id == v
	})
	return nil
}

func (repository *Repository) Create(taskRequest *payloads.TaskRequest) (*payloads.Task, error) {
	var task payloads.Task
	task.ID = repository.getNextID()
	task.Title = taskRequest.Title
	task.Description = taskRequest.Description
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = payloads.StatusCreated
	task.DeleteChan = make(chan struct{})
	repository.DB[task.ID] = &task
	return &task, nil
}

func (repository *Repository) Update(taskRequest *payloads.TaskRequest, id uint64) (*payloads.Task, error) {
	task, ok := repository.DB[id]
	if !ok {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}
	task.Title = taskRequest.Title
	if taskRequest.Description != "" {
		task.Description = taskRequest.Description
	}
	task.UpdatedAt = time.Now()
	return task, nil
}

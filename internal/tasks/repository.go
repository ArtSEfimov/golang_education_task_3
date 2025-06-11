package tasks

import (
	"fmt"
	"io_bound_task/internal/tasks/payloads"
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

func (repository *Repository) GetAllTasks(tasks *payloads.AllTasksResponse) error {
	for _, task := range repository.DB {
		task.SetDuration()
		tasks.Tasks = append(tasks.Tasks, *task)
	}
	return nil
}
func (repository *Repository) GetAllTasksInOrder(tasks *payloads.AllTasksResponse) error {
	for _, id := range repository.Order {
		if task, ok := repository.DB[id]; ok {
			task.SetDuration()
			tasks.Tasks = append(tasks.Tasks, *task)
		}
	}
	return nil
}
func (repository *Repository) Delete(id uint64) error {
	_, ok := repository.DB[id]
	if !ok {
		return fmt.Errorf("task with ID %d not found", id)
	}

	delete(repository.DB, id)
	slices.DeleteFunc(repository.Order, func(v uint64) bool {
		return id == v
	})

	return nil
}

func (repository *Repository) Create(task *payloads.Task) error {
	task.ID = repository.getNextID()
	repository.DB[task.ID] = task
	task.CreatedAt = time.Now()
	task.Status = payloads.StatusCreated
	return nil
}

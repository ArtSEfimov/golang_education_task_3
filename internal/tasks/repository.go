package tasks

import (
	"fmt"
	"slices"
	"time"
)

type MapDB = map[uint64]Task
type OrderedID = []uint64

type Repository struct {
	idCounter uint64
	DB        MapDB
	Order     OrderedID
}

func NewRepository() *Repository {
	db := make(MapDB, 10)
	order := make(OrderedID, 10)
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

func (repository *Repository) FindByID(id uint64) (*Task, error) {
	task, ok := repository.DB[id]
	if !ok {
		return nil, fmt.Errorf("task with id = %d not found", id)
	}
	return &task, nil
}

func (repository *Repository) GetAllTasks(tasks *AllTasksResponse) error {
	for _, task := range repository.DB {
		tasks.Tasks = append(tasks.Tasks, task)
	}
	return nil
}
func (repository *Repository) GetAllTasksInOrder(tasks *AllTasksResponse) error {
	for _, id := range repository.Order {
		if task, ok := repository.DB[id]; ok {
			tasks.Tasks = append(tasks.Tasks, task)
		}
	}
	return nil
}
func (repository *Repository) Delete(id uint64) error {
	_, ok := repository.DB[id]
	if !ok {
		return fmt.Errorf("task with id = %d not found", id)
	}

	delete(repository.DB, id)
	slices.DeleteFunc(repository.Order, func(v uint64) bool {
		return id == v
	})

	return nil
}

func (repository *Repository) Create(task *Task) error {
	task.ID = repository.getNextID()
	repository.DB[task.ID] = *task
	task.CreatedAt = time.Now()
	task.Status = StatusCreated
	return nil
}

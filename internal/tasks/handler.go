package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io_bound_task/internal/tasks/service"
	"io_bound_task/pkg/response"
	"log"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	repository *Repository
	processor  *service.Processor
}

func NewHandlerDeps(repository *Repository, processor *service.Processor) *HandlerDeps {
	return &HandlerDeps{
		repository: repository,
		processor:  processor,
	}
}

type Handler struct {
	repository *Repository
	processor  *service.Processor
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := Handler{
		repository: deps.repository,
		processor:  deps.processor,
	}
	router.HandleFunc("GET /tasks", handler.Get())
	router.HandleFunc("GET /tasks/{id}", handler.GetByID())
	router.HandleFunc("DELETE /tasks/{id}", handler.Delete())
	router.HandleFunc("POST /tasks", handler.Create())
}

// TODO for all GET queries need to update task status and time since start

func (handler *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allTasksResponse AllTasksResponse
		ordered := r.URL.Query().Get("order")
		isOrdered, parseErr := strconv.ParseBool(ordered)
		if parseErr != nil {
			log.Printf("wrong \"order\" parameter value: %s", ordered)
		}

		if parseErr == nil && isOrdered {
			err := handler.repository.GetAllTasksInOrder(&allTasksResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response.JsonResponse(w, allTasksResponse, http.StatusOK)
			return
		}

		err := handler.repository.GetAllTasks(&allTasksResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JsonResponse(w, &allTasksResponse, http.StatusOK)
	}
}

func (handler *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			idErr := fmt.Errorf("wrong id format: %w, get %s", parseErr, idString)
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}

		task, findErr := handler.repository.FindByID(id)
		if findErr != nil {
			http.Error(w, findErr.Error(), http.StatusBadRequest)
			return
		}

		response.JsonResponse(w, task, http.StatusOK)

	}
}

func (handler *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			idErr := fmt.Errorf("wrong id format: %w, get %s", parseErr, idString)
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}

		err := handler.repository.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.JsonResponse(w, nil, http.StatusNoContent)
	}
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&task)
		if decodeErr != nil {
			createErr := fmt.Errorf("bad create query: %w", decodeErr)
			http.Error(w, createErr.Error(), http.StatusBadRequest)
		}
		if task.Name == "" {
			requiredErr := fmt.Errorf("\"name\" value is required")
			http.Error(w, requiredErr.Error(), http.StatusBadRequest)
		}

		creationErr := handler.repository.Create(&task)
		if creationErr != nil {
			http.Error(w, creationErr.Error(), http.StatusInternalServerError)
		}

		handler.processor.AddTask(&task)
		response.JsonResponse(w, &task, http.StatusCreated)

	}
}

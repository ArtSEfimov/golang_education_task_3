package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io_bound_task/internal/tasks/payloads"
	"io_bound_task/internal/tasks/service"
	"io_bound_task/pkg/response"
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
	router.HandleFunc("GET /tasks/{id}", handler.FindByID())
	router.HandleFunc("DELETE /tasks/{id}", handler.Delete())
	router.HandleFunc("POST /tasks", handler.Create())
	router.HandleFunc("PUT /tasks/{id}", handler.Update())
}

func (handler *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validQueryParams := getValidQueryParams(r)
		if isOrdered(r) {
			allTasksResponse, getErr := handler.repository.GetTasksInOrder(validQueryParams)
			if getErr != nil {
				http.Error(w, getErr.Error(), http.StatusInternalServerError)
				return
			}
			response.JsonResponse(w, allTasksResponse, http.StatusOK)
			return
		}

		allTasksResponse, getErr := handler.repository.GetTasks(validQueryParams)
		if getErr != nil {
			http.Error(w, getErr.Error(), http.StatusInternalServerError)
			return
		}

		response.JsonResponse(w, &allTasksResponse, http.StatusOK)
	}
}

func (handler *Handler) FindByID() http.HandlerFunc {
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

		deleteErr := handler.repository.Delete(id, handler.processor)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusBadRequest)
			return
		}

		response.JsonResponse(w, nil, http.StatusNoContent)
	}
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskRequest payloads.TaskRequest
		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&taskRequest)
		if decodeErr != nil {
			createErr := fmt.Errorf("bad create query: %w", decodeErr)
			http.Error(w, createErr.Error(), http.StatusBadRequest)
			return
		}
		if taskRequest.Title == "" {
			requiredErr := fmt.Errorf("\"title\" value is required")
			http.Error(w, requiredErr.Error(), http.StatusBadRequest)
			return
		}

		createdTask, createErr := handler.repository.Create(&taskRequest)
		if createErr != nil {
			http.Error(w, createErr.Error(), http.StatusInternalServerError)
			return
		}

		handler.processor.AddTask(createdTask)

		response.JsonResponse(w, createdTask, http.StatusCreated)
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			idErr := fmt.Errorf("wrong id format: %w, get %s", parseErr, idString)
			http.Error(w, idErr.Error(), http.StatusBadRequest)
			return
		}
		var taskRequest payloads.TaskRequest
		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&taskRequest)
		if decodeErr != nil {
			createErr := fmt.Errorf("bad update query: %w", decodeErr)
			http.Error(w, createErr.Error(), http.StatusBadRequest)
			return
		}

		if taskRequest.Title == "" {
			requiredErr := fmt.Errorf("\"title\" value is required")
			http.Error(w, requiredErr.Error(), http.StatusBadRequest)
			return
		}

		updatedTask, updateErr := handler.repository.Update(&taskRequest, id)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusInternalServerError)
			return
		}

		response.JsonResponse(w, updatedTask, http.StatusOK)
	}
}

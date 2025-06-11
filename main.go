package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io_bound_task/internal/tasks"
	"io_bound_task/internal/tasks/service"
	"net/http"
	"os"
)

func main() {
	loadErr := godotenv.Load(".env")
	if loadErr != nil {
		return
	}

	port := os.Getenv("PORT")

	taskMux := http.NewServeMux()
	taskServer := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: taskMux,
	}

	taskRepository := tasks.NewRepository()
	taskProcessor := service.NewProcessor()
	tasks.NewHandler(
		taskMux, tasks.NewHandlerDeps(
			taskRepository,
			taskProcessor,
		),
	)

	go taskProcessor.Start()

	fmt.Printf("App is listening on port %s...", port)
	listenErr := taskServer.ListenAndServe()
	if listenErr != nil {
		return
	}
}

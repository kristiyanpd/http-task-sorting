package router

import (
	"log"
	"net/http"
	"task/pkg/controller"
	"task/pkg/service"
)

type TaskRouter struct {
	taskController TaskController
}

type TaskController interface {
	SortTasks(http.ResponseWriter, *http.Request)
}

func NewRouter(controller TaskController) *TaskRouter {
	return &TaskRouter{
		taskController: controller,
	}
}

func Init() *TaskRouter {
	taskService := service.NewService()
	taskController := controller.NewController(taskService)
	return NewRouter(taskController)
}

func (r *TaskRouter) Start() {
	// Create router and define routes
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", r.taskController.SortTasks)

	// Start server
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}

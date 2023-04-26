package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"task/pkg/model"
)

type TaskController struct {
	taskService TaskService
}

type TaskService interface {
	SortTasks(tasks []model.Task) ([]*model.Task, error)

	GenerateBashScript(tasks []*model.Task, w io.Writer) error
}

func NewController(service TaskService) *TaskController {
	return &TaskController{
		taskService: service,
	}
}

func (c *TaskController) SortTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed!", http.StatusMethodNotAllowed)
		return
	}

	var request model.TaskRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "JSON input could not be parsed!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	sortedTasks, err := c.taskService.SortTasks(request.Tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Testable using
	// curl -H "Accept: application/bash" -d @examples/tasks.json http://localhost:8000/tasks | bash
	if r.Header.Get("Accept") == "application/bash" {
		w.Header().Set("Content-Type", "application/bash")
		err := c.taskService.GenerateBashScript(sortedTasks, w)
		if err != nil {
			log.Println(err)
		}
		return
	}

	response, err := json.Marshal(sortedTasks)
	if err != nil {
		http.Error(w, "Result could not be serialized to JSON.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
		return
	}
}

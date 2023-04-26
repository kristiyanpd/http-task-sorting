package service

import (
	"errors"
	"fmt"
	"io"
	"task/pkg/model"
)

type TaskService struct{}

func NewService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) SortTasks(tasks []model.Task) ([]*model.Task, error) {
	var sortedTasks = make([]*model.Task, 0, len(tasks))
	var unsortedTasks = make(map[string]*model.Task, len(tasks))

	// Copy initial tasks
	for i := range tasks {
		unsortedTasks[tasks[i].Name] = &tasks[i]
	}

	// Loop until there are no more unsorted tasks
	for len(unsortedTasks) > 0 {
		resolved := false
		for taskName, currentTask := range unsortedTasks {
			if len(currentTask.Requires) == 0 {
				// Place tasks without dependencies in sorted slice and remove from unsorted map
				sortedTasks = append(sortedTasks, currentTask)
				delete(unsortedTasks, taskName)
				// Search for the deleted currentTask in other tasks' dependencies
				for i := range unsortedTasks {
					for j, dep := range unsortedTasks[i].Requires {
						if dep == taskName {
							// Remove deleted currentTask from dependency slice' of other tasks
							unsortedTasks[i].Requires = append(unsortedTasks[i].Requires[:j], unsortedTasks[i].Requires[j+1:]...)
						}
					}
				}
				resolved = true
			}
		}
		// Reachable when all remaining tasks have dependencies
		if !resolved {
			return nil, errors.New("circular dependency found")
		}
	}
	return sortedTasks, nil
}

func (s *TaskService) GenerateBashScript(tasks []*model.Task, w io.Writer) error {
	_, err := fmt.Fprintf(w, "#!/usr/bin/env bash\n")
	if err != nil {
		return err
	}
	for _, task := range tasks {
		_, err := fmt.Fprintf(w, "%s\n", task.Command)
		if err != nil {
			return err
		}
	}
	return nil
}

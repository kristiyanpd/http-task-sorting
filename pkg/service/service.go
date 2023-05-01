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

func (s *TaskService) SortTasksRecursively(tasks []model.Task) ([]*model.Task, error) {
	visited := make(map[string]bool)
	sorted := make([]*model.Task, 0, len(tasks))

	for _, task := range tasks {
		// Pick an unvisited node
		if !visited[task.Name] {
			err := s.topologicalSortHelper(task, tasks, visited, &sorted)
			if err != nil {
				return nil, err
			}
		}
	}
	return sorted, nil
}

func (s *TaskService) topologicalSortHelper(task model.Task, tasks []model.Task, visited map[string]bool, sorted *[]*model.Task) error {
	if visited[task.Name] {
		return errors.New("circular dependency found")
	}
	// Mark task as visited so circular checks can be done
	visited[task.Name] = true

	// Go through current task dependencies
	for _, req := range task.Requires {
		// Check if task is already sorted
		if !s.containsTask(req, *sorted) {
			reqTask := s.findTaskByName(req, tasks)
			// DFS on required task
			err := s.topologicalSortHelper(*reqTask, tasks, visited, sorted)
			if err != nil {
				return err
			}
		}
	}

	// Add task to sorted tasks so when iterating through other tasks' dependencies, we can skip it
	*sorted = append(*sorted, &task)
	return nil
}

func (s *TaskService) findTaskByName(taskName string, tasks []model.Task) *model.Task {
	for _, task := range tasks {
		if task.Name == taskName {
			return &task
		}
	}
	return &model.Task{}
}

func (s *TaskService) containsTask(taskName string, tasks []*model.Task) bool {
	for _, task := range tasks {
		if task.Name == taskName {
			return true
		}
	}
	return false
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

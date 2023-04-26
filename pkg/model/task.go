package model

type Task struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Requires []string `json:"requires,omitempty"`
}

type TaskRequest struct {
	Tasks []Task `json:"tasks"`
}

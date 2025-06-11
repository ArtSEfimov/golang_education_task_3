package tasks

type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AllTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

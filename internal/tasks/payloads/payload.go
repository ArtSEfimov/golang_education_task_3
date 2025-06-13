package payloads

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AllTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

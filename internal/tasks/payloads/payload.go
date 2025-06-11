package payloads

type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AllTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

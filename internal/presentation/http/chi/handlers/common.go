package handlers

type Status string

const (
	Success Status = "success"
	Failed  Status = "failed"
)

type Response struct {
	Data   any      `json:"data,omitempty"`
	Status Status   `json:"status"`
	Errors []string `json:"errors,omitempty"`
}

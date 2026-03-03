package model

type ContainerState struct {
	Status     string `json:"status"`
	Running    bool   `json:"running"`
	Paused     bool   `json:"paused"`
	Restarting bool   `json:"restarting"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

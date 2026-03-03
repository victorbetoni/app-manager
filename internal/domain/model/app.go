package model

type App struct {
	ID          string `json:"id"`
	ContainerID string `json:"container_id"`
	Branch      string `json:"branch"`
	Name        string `json:"name"`
	Dir         string `json:"dir"`
}

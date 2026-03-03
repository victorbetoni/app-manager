package model

type Request struct {
	AppID    string `json:"app"`
	Action   string `json:"action"`
	TailSize *int   `json:"tail_size"`
}

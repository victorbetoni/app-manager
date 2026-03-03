package model

type Message struct {
	Error     bool   `json:"error"`
	Info      bool   `json:"info"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

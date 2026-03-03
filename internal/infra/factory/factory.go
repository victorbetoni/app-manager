package factory

import (
	"github.com/victorbetoni/justore/app-manager/internal/infra/action"
)

func CreateProcessRequestStrategy(act string) action.ProcessActionStrategy {
	switch act {
	case "logs":
		return action.ProcessLogs{}
	default:
		return nil
	}
}

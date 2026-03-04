package action

import (
	"context"
	"fmt"
	"time"

	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
	"github.com/victorbetoni/justore/app-manager/internal/infra/command"
)

type ProcessLogs struct{}

func (p ProcessLogs) Process(ctx context.Context, output chan model.Message, app model.App, tailSize int) {

	cmd := "docker"
	args := []string{"logs", "--tail", fmt.Sprintf("%d", tailSize), "-f", app.ContainerID}

	c := command.NewCommand(cmd, args)
	ch := make(chan []byte)
	if err := c.StreamExec(ctx, ch); err != nil {
		output <- model.Message{
			Error:     true,
			Timestamp: time.Now().Unix(),
			Message:   err.Error(),
		}
		return
	}
	for data := range ch {
		msg := model.Message{
			Error:     false,
			Timestamp: time.Now().Unix(),
			Message:   string(data),
		}

		select {
		case <-ctx.Done():
			return
		case output <- msg:
		}
	}

}

func (p ProcessLogs) Streamable() bool {
	return true
}

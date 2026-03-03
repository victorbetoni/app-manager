package docker

import (
	"context"

	"github.com/moby/moby/client"
	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
)

var Client *client.Client

func init() {
	cli, err := client.New()
	if err != nil {
		panic(err)
	}
	Client = cli
}

func ContainerState(ctx context.Context, id string) (*model.ContainerState, error) {

	c, err := Client.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return nil, err
	}

	return &model.ContainerState{
		Status:     string(c.Container.State.Status),
		Running:    c.Container.State.Running,
		Paused:     c.Container.State.Paused,
		Restarting: c.Container.State.Restarting,
		StartedAt:  c.Container.State.StartedAt,
		FinishedAt: c.Container.State.FinishedAt,
	}, nil

}

func ContainerStart(ctx context.Context, id string) error {
	_, err := Client.ContainerStart(ctx, id, client.ContainerStartOptions{})
	return err
}

func ContainerRestart(ctx context.Context, id string) error {
	_, err := Client.ContainerRestart(ctx, id, client.ContainerRestartOptions{})
	return err
}

func ContainerStop(ctx context.Context, id string) error {
	_, err := Client.ContainerStop(ctx, id, client.ContainerStopOptions{})
	return err
}

package action

import (
	"context"

	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
)

type ProcessActionStrategy interface {
	Process(ctx context.Context, output chan model.Message, app model.App, tailSize int)
	Streamable() bool
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
	"github.com/victorbetoni/justore/app-manager/internal/infra/app"
	"github.com/victorbetoni/justore/app-manager/internal/infra/docker"
)

func ListApps(ctx *gin.Context) {

	type AppList struct {
		App   model.App             `json:"app"`
		State *model.ContainerState `json:"state"`
	}

	type Grid struct {
		Data       []AppList `json:"data"`
		TotalCount int       `json:"total_count"`
	}

	apps := make([]AppList, 0)

	for _, a := range app.Apps {
		state, err := docker.ContainerState(ctx, a.ContainerID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		apps = append(apps, AppList{
			App:   a,
			State: state,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "body": Grid{
		Data:       apps,
		TotalCount: len(apps),
	}})

}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
	"github.com/victorbetoni/justore/app-manager/internal/infra/app"
	"github.com/victorbetoni/justore/app-manager/internal/infra/docker"
)

func ContainerState(ctx *gin.Context) {

	appId := ctx.Param("appId")

	type AppList struct {
		App   model.App             `json:"app"`
		State *model.ContainerState `json:"state"`
	}

	app, ok := app.Apps[appId]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "App não encontrado"})
		return
	}

	state, err := docker.ContainerState(ctx, app.ContainerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "body": AppList{
		App:   app,
		State: state,
	}})

}

func ContainerStop(ctx *gin.Context) {

	appId := ctx.Param("appId")

	app, ok := app.Apps[appId]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "App não encontrado"})
		return
	}

	if err := docker.ContainerStop(ctx, app.ContainerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Container parado com sucesso"})

}

func ContainerStart(ctx *gin.Context) {

	appId := ctx.Param("appId")

	app, ok := app.Apps[appId]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "App não encontrado"})
		return
	}

	if err := docker.ContainerStart(ctx, app.ContainerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Container iniciado com sucesso"})

}

func ContainerRestart(ctx *gin.Context) {

	appId := ctx.Param("appId")

	app, ok := app.Apps[appId]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "App não encontrado"})
		return
	}

	if err := docker.ContainerRestart(ctx, app.ContainerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Container reiniciado com sucesso"})

}

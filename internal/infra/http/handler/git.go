package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorbetoni/justore/app-manager/internal/infra/app"
	"github.com/victorbetoni/justore/app-manager/internal/infra/command"
)

func SynchronizeProject(ctx *gin.Context) {

	appId := ctx.Param("appId")

	app, ok := app.Apps[appId]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "App não encontrado"})
		return
	}

	cmdReset := command.NewDirCommand("git", []string{"reset", "--hard", "HEAD"}, app.Dir)
	if _, err := cmdReset.Exec(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	cmd := "git"
	args := []string{"pull", "origin", app.Branch}

	c := command.NewDirCommand(cmd, args, app.Dir)
	if _, err := c.Exec(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Projeto sincronizado com sucesso"})
}

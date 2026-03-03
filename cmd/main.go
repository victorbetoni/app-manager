package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/victorbetoni/justore/app-manager/internal/config"
	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
	"github.com/victorbetoni/justore/app-manager/internal/infra/app"
	"github.com/victorbetoni/justore/app-manager/internal/infra/docker"
	"github.com/victorbetoni/justore/app-manager/internal/infra/http"
	"github.com/victorbetoni/justore/app-manager/internal/infra/jwt"
	"github.com/victorbetoni/justore/app-manager/internal/infra/ws"
)

func main() {

	config.Load()

	f, err := os.ReadFile(config.GetConfig().AppsJsonPath)
	if err != nil {
		panic(err)
	}

	var apps []model.App
	if err := json.Unmarshal(f, &apps); err != nil {
		panic(err)
	}

	for _, a := range apps {
		app.Apps[a.ID] = a
	}

	if config.GetConfig().UseAuth {
		pubKeyFile, err := os.ReadFile(config.GetConfig().Jwt.Keys.PublicKey)
		if err != nil {
			log.Fatalf("Couldnt retrieve RSA keys: %s\n", err.Error())
			return
		}
		privateKeyFile, err := os.ReadFile(config.GetConfig().Jwt.Keys.PrivateKey)
		if err != nil {
			log.Fatalf("Couldnt retrieve RSA keys: %s\n", err.Error())
			return
		}
		jwt.DefineKeyPair(privateKeyFile, pubKeyFile)
	}

	defer func() {
		if docker.Client != nil {
			docker.Client.Close()
		}
	}()

	connectionHub := ws.NewHub()

	r := http.Build(connectionHub)

	if config.GetConfig().UseTLS {
		err = r.RunTLS(fmt.Sprintf(":%d", config.GetConfig().Port), config.GetConfig().TlsCert.PublicKey, config.GetConfig().TlsCert.PrivateKey)
	} else {
		err = r.Run(fmt.Sprintf(":%d", config.GetConfig().Port))
	}
	if err != nil {
		log.Fatalf("Couldn't start HTTP server: %s\n", err.Error())
		return
	}

}

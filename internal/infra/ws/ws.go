package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/gorilla/websocket"
	"github.com/victorbetoni/justore/app-manager/internal/domain/model"
	"github.com/victorbetoni/justore/app-manager/internal/infra/app"
	"github.com/victorbetoni/justore/app-manager/internal/infra/factory"
)

type WebsocketSession struct {
	Context       context.Context
	Identifier    string
	IP            string
	UserAgent     string
	Chan          chan model.Message
	Conn          *websocket.Conn
	connectionHub *ConnectionHub
	LastActivity  time.Time
}

func NewSession(ctx context.Context, identifier, ip, userAgent string, conn *websocket.Conn, h *ConnectionHub) *WebsocketSession {
	return &WebsocketSession{
		Context:       ctx,
		Identifier:    identifier,
		LastActivity:  time.Now(),
		Chan:          make(chan model.Message),
		IP:            ip,
		UserAgent:     userAgent,
		Conn:          conn,
		connectionHub: h,
	}
}

func (s *WebsocketSession) Listen() error {
	defer s.connectionHub.Unregister(s.Conn)
	defer close(s.Chan)

	go func() {
		for msg := range s.Chan {
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error while parsing message")
				continue
			}
			s.Conn.WriteMessage(1, data)
		}
	}()

	for {

		_, p, err := s.Conn.ReadMessage()
		if err != nil {
			return err
		}

		var req model.Request
		if err := json.Unmarshal(p, &req); err != nil {
			fmt.Println(err.Error())
		}

		a, ok := app.Apps[req.AppID]
		if !ok {
			continue
		}

		strategy := factory.CreateProcessRequestStrategy(req.Action)
		if strategy == nil {
			continue
		}

		if req.TailSize != nil {
			*req.TailSize = int(float64(math.Min(math.Max(float64(*req.TailSize), 0), 10000)))
		}

		if strategy.Streamable() {
			go strategy.Process(s.Context, s.Chan, a, *req.TailSize) // se for streamable roda em uma routine separada e nao espera pelo fim do comando
		} else {
			strategy.Process(s.Context, s.Chan, a, *req.TailSize)
		}

	}
}

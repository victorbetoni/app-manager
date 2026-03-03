package ws

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionHub struct {
	clients map[*websocket.Conn]*WebsocketSession
	mu      sync.Mutex
}

func NewHub() *ConnectionHub {
	return &ConnectionHub{
		clients: make(map[*websocket.Conn]*WebsocketSession),
	}
}

func (h *ConnectionHub) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		_ = client.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *ConnectionHub) Register(ctx context.Context, dcId string, ip string, userAgent string, conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = NewSession(ctx, dcId, ip, userAgent, conn, h)
	h.mu.Unlock()
}

func (h *ConnectionHub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	if _, ok := h.clients[conn]; ok {
		delete(h.clients, conn)
		conn.Close()
	}
	h.mu.Unlock()
}

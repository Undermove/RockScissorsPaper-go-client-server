package auth

import (
	"../data"
	"github.com/gorilla/websocket"
)

type AuthModule struct {
	Clients     map[*websocket.Conn]string
	AuthClients map[string]*websocket.Conn
}

func NewModule() *AuthModule {
	return &AuthModule{
		Clients:     make(map[*websocket.Conn]string),
		AuthClients: make(map[string]*websocket.Conn),
	}
}

func (a *AuthModule) AddConnection(w *websocket.Conn) bool {
	if _, ok := a.Clients[w]; ok {
		return false
	}
	a.Clients[w] = "unknown"
	return true
}

func (a *AuthModule) Disconnect(w *websocket.Conn) {
	username := a.Clients[w]
	delete(a.Clients, w)
	delete(a.AuthClients, username)
}

func (a *AuthModule) IsLoggedIn(w *websocket.Conn) bool {
	if client, ok := a.Clients[w]; ok {
		if client != "unknown" {
			return true
		}
		return false
	}

	return false
}

func (a *AuthModule) Authenticate(w *websocket.Conn, m msg.AuthRequest) bool {
	if _, ok := a.AuthClients[m.Username]; ok {
		return false
	}

	a.AuthClients[m.Username] = w
	a.Clients[w] = m.Username

	return true
}

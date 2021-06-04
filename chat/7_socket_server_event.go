package chat

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ravielze/otopal/auth"
)

func (server *ChatServer) OnDisconnect(conn *websocket.Conn) {
	server.Lock()
	defer server.Unlock()
	for i, c := range server.connection {
		if c.connection == conn {
			server.module.controller.OnDisconnect(c)
			delete(server.connection, i)
		}
	}
}

func (server *ChatServer) OnConnect(conn *websocket.Conn, user auth.User, exp int64) *SocketConnection {
	server.Lock()
	defer server.Unlock()
	so := NewConnection(conn, server.lastId, user, exp)
	server.connection[server.lastId] = &so
	server.Broadcast(struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("%d connected.", server.lastId),
	})
	server.lastId++
	server.Refresh(conn)
	server.module.controller.OnConnect(&so)
	return &so
}

func (server *ChatServer) Refresh(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	conn.SetWriteDeadline(time.Now().Add(time.Minute * 3))
}

func (server *ChatServer) Broadcast(msg interface{}) {
	for _, c := range server.connection {
		c.connection.WriteJSON(msg)
	}
}

func (server *ChatServer) GetConnection(socketId int) *SocketConnection {
	server.Lock()
	defer server.Unlock()
	return server.connection[socketId]
}

func (server *ChatServer) GetConnectionByUser(userId uint) []*SocketConnection {
	server.Lock()
	defer server.Unlock()
	var result []*SocketConnection
	for _, c := range server.connection {
		if c.User.ID == userId {
			result = append(result, c)
		}
	}
	return result
}

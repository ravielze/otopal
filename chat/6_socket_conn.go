package chat

import (
	"github.com/gorilla/websocket"
	"github.com/ravielze/otopal/auth"
)

type (
	SocketConnection struct {
		connection *websocket.Conn
		server     *ChatServer

		ID   int
		User auth.User
	}
)

func NewConnection(conn *websocket.Conn, id int, user auth.User) SocketConnection {
	return SocketConnection{
		connection: conn,
		ID:         id,
		User:       user,
		server:     ChatServerInstance,
	}
}

func (s *SocketConnection) Message(data interface{}) {
	s.connection.WriteJSON(data)
}

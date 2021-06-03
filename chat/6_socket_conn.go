package chat

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/ravielze/otopal/auth"
)

type (
	SocketConnection struct {
		connection *websocket.Conn
		server     *ChatServer

		ID      int
		User    auth.User
		Expired int64
	}
)

func NewConnection(conn *websocket.Conn, id int, user auth.User, exp int64) SocketConnection {
	return SocketConnection{
		connection: conn,
		ID:         id,
		User:       user,
		server:     ChatServerInstance,
		Expired:    exp,
	}
}

func (s *SocketConnection) Message(data interface{}) {
	s.connection.WriteJSON(data)
}

func (s *SocketConnection) IsValid() bool {
	return time.Now().Unix() < s.Expired
}

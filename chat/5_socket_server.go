package chat

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/otopal/auth"

	"github.com/gorilla/websocket"
)

type ChatServer struct {
	sync.RWMutex
	Running    chan os.Signal
	module     Module
	connection map[int]*SocketConnection
	lastId     int
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var ChatServerInstance *ChatServer

func NewChatServer() *ChatServer {
	result := &ChatServer{
		Running:    make(chan os.Signal, 1),
		module:     module_manager.GetModule("chat").(Module),
		connection: make(map[int]*SocketConnection),
		lastId:     1,
	}
	return result
}

func (server *ChatServer) Run(g *gin.Engine) {
	auc := module_manager.GetModule("auth").(auth.Module).Usecase()
	g.GET("/chat",
		auc.AuthenticationNeeded(),
		auc.AllowedRole(auth.ROLE_TECHNICIAN, auth.ROLE_CUSTOMER),
		func(c *gin.Context) {
			user := auc.GetUser(c)
			exp := c.Keys["exp"].(int64)
			go server.websocketHandler(c.Writer, c.Request, user, exp)
		})
}

func (server *ChatServer) Handle(payload *StandardPayload, socketConn *SocketConnection) {
	if !socketConn.IsValid() {
		server.OnDisconnect(socketConn.connection)
		socketConn.connection.Close()
		return
	}
	switch strings.ToLower(payload.Event) {
	case "send":
		server.module.controller.OnSendMessage(socketConn, payload.Data)
	case "read":
		server.module.controller.OnReadMessage(socketConn, payload.Data)
	}
}

func (server *ChatServer) websocketHandler(w http.ResponseWriter, r *http.Request, user auth.User, exp int64) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	so := server.OnConnect(conn, user, exp)
	var payload StandardPayload
	for {
		payload = StandardPayload{}
		err := conn.ReadJSON(&payload)
		if err != nil {
			break
		}
		server.Handle(&payload, so)
		server.Refresh(conn)
	}
	server.OnDisconnect(conn)
}

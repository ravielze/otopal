package chat

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	connection map[int]SocketConnection
	lastId     int
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var ChatServerInstance *ChatServer

func NewChatServer() *ChatServer {
	result := &ChatServer{
		Running:    make(chan os.Signal, 1),
		module:     module_manager.GetModule("chat").(Module),
		connection: make(map[int]SocketConnection),
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
			go server.websocketHandler(c.Writer, c.Request, user)
		})
}

func (server *ChatServer) websocketHandler(w http.ResponseWriter, r *http.Request, user auth.User) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	server.OnConnect(conn, user)
	// conn.SetCloseHandler(func(code int, text string) error {
	// 	for _, c := range server.connection {
	// 		c.WriteMessage(t, []byte(response))
	// 	}
	// })
	for {
		var payload StandardPayload
		err := conn.ReadJSON(&payload)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			break
		}
		fmt.Println(payload)
		server.Refresh(conn)
	}
	server.OnDisconnect(conn)
}

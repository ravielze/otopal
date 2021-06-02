package chat

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	module_manager "github.com/ravielze/oculi/common/module"

	"github.com/gorilla/websocket"
)

type ChatServer struct {
	sync.RWMutex
	Running    chan os.Signal
	module     Module
	connection map[int]*websocket.Conn
	lastId     int
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewChatServer() *ChatServer {
	result := &ChatServer{
		Running:    make(chan os.Signal, 1),
		module:     module_manager.GetModule("chat").(Module),
		connection: make(map[int]*websocket.Conn),
		lastId:     1,
	}
	return result
}

func (server *ChatServer) Run(g *gin.Engine) {
	g.GET("/chat", func(c *gin.Context) {
		go server.websocketHandler(c.Writer, c.Request)
	})
}

func (server *ChatServer) OnDisconnect(conn *websocket.Conn) {
	server.Lock()
	defer server.Unlock()
	for i, c := range server.connection {
		if c == conn {
			delete(server.connection, i)
		}
	}
}

func (server *ChatServer) OnConnect(conn *websocket.Conn) {
	server.Lock()
	defer server.Unlock()
	server.connection[server.lastId] = conn
	server.Broadcast(struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("%d connected.", server.lastId),
	})
	server.lastId++
	server.Refresh(conn)
}

func (server *ChatServer) Refresh(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
}

func (server *ChatServer) Broadcast(msg interface{}) {
	for _, c := range server.connection {
		c.WriteJSON(msg)
	}
}

func (server *ChatServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	server.OnConnect(conn)
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

type (
	StandardPayload struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}
)

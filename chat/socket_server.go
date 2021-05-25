package chat

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ChatServer struct {
	server  *socketio.Server
	Running chan os.Signal
}

var server *socketio.Server

func NewChatServer() *ChatServer {
	server = socketio.NewServer(nil)
	eventHandler(server)
	return &ChatServer{
		server:  server,
		Running: make(chan os.Signal, 1),
	}
}

func ginMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func (cs *ChatServer) Run(g *gin.Engine, allowedAddress []string) {
	socketGroup := g.Group("/socket.io")
	socketGroup.Use(ginMiddleware(allowedAddress[0]))
	socketGroup.GET("/*any", gin.WrapH(cs.server))
	socketGroup.POST("/*any", gin.WrapH(cs.server))
	go func() {
		if err := cs.server.Serve(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Chat Server running...")
	defer cs.server.Close()
	<-cs.Running
}

func eventHandler(server *socketio.Server) {

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		fmt.Println(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnEvent("/", "gotoChat", func(s socketio.Conn) {
		s.Join("/chat")
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
}

package chat

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type Controller struct {
	uc IUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc: uc,
	}
	return cont
}

func (cont Controller) OnConnect(s socketio.Conn) error {
	panic("not implemented")
}

func (cont Controller) OnDisconnect(s socketio.Conn) error {
	panic("not implemented")
}

func (cont Controller) OnReadMessage(s socketio.Conn, msg string) string {
	panic("not implemented")
}

func (cont Controller) OnRetrieveMessage(s socketio.Conn, msg string) string {
	panic("not implemented")
}

func (cont Controller) OnSendMessage(s socketio.Conn, msg string) string {
	panic("not implemented")
}

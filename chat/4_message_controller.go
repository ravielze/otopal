package chat

import (
	"fmt"

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
	fmt.Printf("User %s connected to chat server.", s.ID())
	return nil
}

func (cont Controller) OnDisconnect(s socketio.Conn) error {
	user, err := cont.uc.GetUserID(s.ID())
	if err == nil {
		cont.uc.Logout(user, s.ID())
	}
	return nil
}

func (cont Controller) OnReadMessage(s socketio.Conn, msg string) string {
	panic("not implemented")
}

func (cont Controller) OnRetrieveMessage(s socketio.Conn, msg string) string {
	panic("not implemented")
}

func (cont Controller) OnSendMessage(s socketio.Conn, msg string) string {
	//binding.JSON.BindBody([]byte(msg))
}

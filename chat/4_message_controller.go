package chat

import (
	"github.com/gin-gonic/gin"
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

// func (cont Controller) OnConnect(s socketio.Conn) error {
// 	fmt.Printf("User %s connected to chat server.", s.ID())
// 	return nil
// }

// func (cont Controller) OnDisconnect(s socketio.Conn, reason string) {
// 	user, err := cont.uc.GetUserID(s.ID())
// 	if err == nil {
// 		cont.uc.Logout(user, s.ID())
// 	}
// }

// func (cont Controller) OnReadMessage(s socketio.Conn, msg string) string {
// 	panic("not implemented")
// }

// func (cont Controller) OnRetrieveMessage(s socketio.Conn, msg string) string {
// 	panic("not implemented")
// }

// func (cont Controller) OnSendMessage(s socketio.Conn, msg string) string {
// 	//binding.JSON.BindBody([]byte(msg))
// 	return ""
// }

// func (cont Controller) OnLogin(s socketio.Conn, msg string) string {
// 	return ""
// }

// func (cont Controller) OnLogout(s socketio.Conn, msg string) string {
// 	return ""
// }

package chat

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/otopal/auth"
)

type Controller struct {
	uc  IUsecase
	auc auth.IUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc:  uc,
		auc: module_manager.GetModule("auth").(auth.Module).Usecase(),
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

func (cont Controller) OnConnect(s *SocketConnection) {
	cont.uc.Login(s.User.ID, s.ID)
	s.server.Broadcast(
		NewStatusPayload(&s.User, true),
	)
}

func (cont Controller) OnDisconnect(s *SocketConnection) {
	cont.uc.Logout(s.User.ID, s.ID)
	s.server.Broadcast(
		NewStatusPayload(&s.User, false),
	)
}

func (cont Controller) OnReadMessage(s *SocketConnection, data interface{}) {
	buff, err := json.Marshal(data)
	if err != nil {
		s.Message(NewErrorPayload(&s.User, err))
		return
	}
	var payload ReadMessageRequestPayload
	err = json.Unmarshal(buff, &payload)
	if err != nil {
		s.Message(NewErrorPayload(&s.User, err))
		return
	}

	sender, err2 := cont.auc.GetRawUser(payload.SenderID)
	if err2 != nil {
		s.Message(NewErrorPayload(&s.User, err2))
		return
	}

	err = cont.uc.ReadAll(s.User.ID, payload.SenderID)
	if err != nil {
		s.Message(NewErrorPayload(&s.User, err))
		return
	}

	s.Message(NewReadPayload(&s.User, &sender))
	if cont.uc.IsOnline(sender.ID) {
		for _, c := range s.server.GetConnectionByUser(sender.ID) {
			if c == nil {
				continue
			}
			c.Message(NewReadPayload(&s.User, &sender))
		}
	}
}

func (cont Controller) RetrieveMessage(ctx *gin.Context) {
	panic("not implemented")
}

func (cont Controller) Overview(ctx *gin.Context) {
	panic("a")
}

func (cont Controller) OnSendMessage(s *SocketConnection, data interface{}) {
	panic("not implemented")
}

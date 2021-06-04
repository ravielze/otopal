package chat

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common/controller_utils"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/oculi/common/utils"
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
	chatGroup := g.Group("/message")
	chatGroup.Use(cont.auc.AuthenticationNeeded(), cont.auc.AllowedRole(auth.ROLE_CUSTOMER, auth.ROLE_TECHNICIAN))
	{
		chatGroup.GET("/overview", cont.Overview)
		chatGroup.GET("/retrieve/:otherUserId", cont.RetrieveMessage)
	}
	return cont
}

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
	ok, params, _ := controller_utils.NewControlChain(ctx).Param("otherUserId").End()
	if ok {
		otherUserId, err := strconv.Atoi(params["otherUserId"])
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		user := cont.auc.GetUser(ctx)
		result, err2 := cont.uc.GetMessage(user.ID, uint(otherUserId))
		if err2 != nil {
			utils.AbortUsecaseError(ctx, err2)
			return
		}
		otherUser, err3 := cont.auc.GetRawUser(uint(otherUserId))
		if err3 != nil {
			utils.AbortUsecaseError(ctx, err3)
			return
		}
		utils.OKAndResponseData(ctx, NewMessageRetrieveResponse(result, user, otherUser))
	}
}

func (cont Controller) Overview(ctx *gin.Context) {
	user := cont.auc.GetUser(ctx)
	resultMsg, resultUnread, err2 := cont.uc.GetOverview(user.ID)
	if err2 != nil {
		utils.AbortUsecaseError(ctx, err2)
		return
	}

	overview := make([]SubMessageOverviewResponse, len(resultMsg))
	for i := range overview {
		overview[i] = resultMsg[i].ConvertOverview(user.ID, resultUnread[i])
	}

	utils.OKAndResponseData(ctx, NewMessageOverviewResponse(overview, user))
}

func (cont Controller) OnSendMessage(s *SocketConnection, data interface{}) {
	buff, err := json.Marshal(data)
	if err != nil {
		s.Message(NewErrorPayload(&s.User, err))
		return
	}
	var payload SendMessageRequestPayload
	err = json.Unmarshal(buff, &payload)
	if err != nil {
		s.Message(NewErrorPayload(&s.User, err))
		return
	}

	receiver, err2 := cont.auc.GetRawUser(payload.ReceiverID)
	if err2 != nil {
		s.Message(NewErrorPayload(&s.User, err2))
		return
	}

	msg, err3 := cont.uc.SendMessage(s.User.ID, payload.ReceiverID, payload.Message)
	if err3 != nil {
		s.Message(NewErrorPayload(&s.User, err3))
		return
	}

	s.Message(NewSendPayload(&s.User, &receiver, msg.Message, msg.CreatedAt.Format("02-01-2006 15:04:05")))
	if cont.uc.IsOnline(receiver.ID) {
		for _, c := range s.server.GetConnectionByUser(receiver.ID) {
			if c == nil {
				continue
			}
			c.Message(NewSendPayload(&s.User, &receiver, msg.Message, msg.CreatedAt.Format("02-01-2006 15:04:05")))
		}
	}
}

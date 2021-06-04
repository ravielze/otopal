package chat

import "github.com/ravielze/otopal/auth"

type (
	StandardPayload struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}

	MessageResponse struct {
		Message string `json:"message"`
		Sent    string `json:"sent"`
		Mine    bool   `json:"mine"`
	}

	MessagesRetrieveResponse struct {
		Sender   auth.UserResponse `json:"sender"`
		Receiver auth.UserResponse `json:"receiver"`
		Messages []MessageResponse `json:"messages"`
	}

	MessageOverviewResponse struct {
		Myself   auth.UserResponse            `json:"myself"`
		Overview []SubMessageOverviewResponse `json:"overview"`
	}

	SubMessageOverviewResponse struct {
		User        auth.UserResponse `json:"user"`
		LastMessage string            `json:"last_message"`
		LastSent    string            `json:"last_sent"`
		Unread      uint              `json:"unread"`
	}

	ResponsePayload struct {
		Payload int               `json:"payload_id"`
		Sender  auth.UserResponse `json:"user"`
	}

	SendMessageRequestPayload struct {
		Message    string `json:"message"`
		ReceiverID uint   `json:"receiver_id"`
	}

	SendMessageResponsePayload struct {
		ResponsePayload
		Receiver auth.UserResponse `json:"receiver"`
		Message  string            `json:"message"`
		Sent     string            `json:"sent"`
	}

	ReadMessageRequestPayload struct {
		SenderID uint `json:"sender_id"`
	}

	ReadMessageResponsePayload struct {
		ResponsePayload
		Receiver auth.UserResponse `json:"receiver"`
	}

	StatusResponsePayload struct {
		ResponsePayload
		Status bool `json:"is_online"`
	}

	ErrorResponsePayload struct {
		ResponsePayload
		Detail string `json:"error_detail"`
	}
)

const (
	ERROR  int = 0
	STATUS int = 1
	READ   int = 2
	SEND   int = 3
)

func (item Message) Convert(senderId uint) MessageResponse {
	mine := false
	if item.UserID == senderId {
		mine = true
	}
	return MessageResponse{
		Message: item.Message,
		Sent:    item.CreatedAt.Format("02-01-2006 15:04:05"),
		Mine:    mine,
	}
}

func (item Message) ConvertOverview(senderId uint, unread uint) SubMessageOverviewResponse {
	var user auth.User
	if item.UserID == senderId {
		user = item.Receiver
	} else {
		user = item.User
	}
	return SubMessageOverviewResponse{
		User:        user.Convert(),
		LastMessage: item.Message,
		LastSent:    item.CreatedAt.Format("02-01-2006 15:04:05"),
		Unread:      unread,
	}
}

func NewMessageOverviewResponse(item []SubMessageOverviewResponse, sender auth.User) MessageOverviewResponse {
	return MessageOverviewResponse{
		Myself:   sender.Convert(),
		Overview: item,
	}
}

func NewMessageRetrieveResponse(item []Message, sender, receiver auth.User) MessagesRetrieveResponse {
	if len(item) == 0 {
		return MessagesRetrieveResponse{
			Sender:   sender.Convert(),
			Receiver: receiver.Convert(),
			Messages: make([]MessageResponse, 0),
		}
	}

	msg := make([]MessageResponse, len(item))
	for i, x := range item {
		msg[i] = x.Convert(sender.ID)
	}

	return MessagesRetrieveResponse{
		Sender:   sender.Convert(),
		Receiver: receiver.Convert(),
		Messages: msg,
	}
}

func NewErrorPayload(user *auth.User, err error) ErrorResponsePayload {
	return ErrorResponsePayload{
		ResponsePayload: ResponsePayload{
			Payload: ERROR,
			Sender:  user.Convert(),
		},
		Detail: err.Error(),
	}
}

func NewStatusPayload(user *auth.User, status bool) StatusResponsePayload {
	return StatusResponsePayload{
		ResponsePayload: ResponsePayload{
			Payload: STATUS,
			Sender:  user.Convert(),
		},
		Status: status,
	}
}

func NewReadPayload(sender, receiver *auth.User) ReadMessageResponsePayload {
	return ReadMessageResponsePayload{
		ResponsePayload: ResponsePayload{
			Payload: READ,
			Sender:  sender.Convert(),
		},
		Receiver: receiver.Convert(),
	}
}

func NewSendPayload(sender, receiver *auth.User, msg, sent string) SendMessageResponsePayload {
	return SendMessageResponsePayload{
		ResponsePayload: ResponsePayload{
			Payload: SEND,
			Sender:  sender.Convert(),
		},
		Receiver: receiver.Convert(),
		Message:  msg,
		Sent:     sent,
	}
}

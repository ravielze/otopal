package chat

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/auth"
)

type Message struct {
	common.UUIDBase `gorm:"embedded;embeddedPrefix:blog_"`
	CreatedAt       time.Time `json:"created_at"`
	Message         string    `gorm:"type:VARCHAR(2097152);"`
	UserID          uint
	User            auth.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReceiverID      uint
	Receiver        auth.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Read            bool
}

func (Message) TableName() string {
	return "message"
}

type IController interface {
	RetrieveMessage(ctx *gin.Context)
	Overview(ctx *gin.Context)

	OnConnect(s *SocketConnection)
	OnDisconnect(s *SocketConnection)
	OnSendMessage(s *SocketConnection, data interface{})
	OnReadMessage(s *SocketConnection, data interface{})
}

type IUsecase interface {
	SendMessage(userId uint, receiverId uint, message string) (Message, error)
	ReadAll(userId uint, receiverId uint) error

	Login(userId uint, socketId int) error
	Logout(userId uint, socketId int) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)
	GetOverview(userId uint) ([]Message, []uint, error)

	IsOnline(userId uint) bool

	GetUserID(socketId int) (uint, error)
}

type IRepo interface {
	GetUserID(socketId int) (uint, error)
	IsLoggedIn(userId uint) bool
	Online(userId uint, socketId int)
	Offline(userId uint, socketId int)

	CreateMessage(msg Message) (Message, error)
	ReadAll(receiverId uint, senderId uint) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)
	//GetOverview(userId uint) ([]uint, error)
	GetUnreadMessage(senderId, receiverId uint) (uint, error)
	GetSender(userId uint) ([]uint, error)
	GetLastMessage(userId uint, user2Id uint) (Message, error)
}

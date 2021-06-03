package chat

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	OnConnect(s *websocket.Conn) error
	OnDisconnect(s *websocket.Conn)

	OnRetrieveMessage(ctx *gin.Context)

	OnSendMessage(s *websocket.Conn, data interface{}) string
	OnReadMessage(s *websocket.Conn, data interface{}) string
	OnLogin(s *websocket.Conn, data interface{}) string
	OnLogout(s *websocket.Conn, data interface{}) string
}

type IUsecase interface {
	SendMessage(userId uint, receiverId uint, message string) (Message, error)
	ReadAll(userId uint, receiverId uint) error

	Login(userId uint, socketId int) error
	Logout(userId uint, socketId int) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)

	IsOnline(userId uint) bool

	GetUserID(socketId int) (uint, error)
}

type IRepo interface {
	GetUserID(socketId int) (uint, error)
	IsLoggedIn(userId uint) bool
	Online(userId uint, socketId int)
	Offline(userId uint, socketId int)

	CreateMessage(msg Message) (Message, error)
	ReadAll(userId uint, senderId uint) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)
	//GetOverview(userId uint) ([]uint, error)
}

package chat

import (
	"time"

	socketio "github.com/googollee/go-socket.io"
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
	OnConnect(s socketio.Conn) error
	OnDisconnect(s socketio.Conn) error
	OnRetrieveMessage(s socketio.Conn, msg string) string
	OnSendMessage(s socketio.Conn, msg string) string
	OnReadMessage(s socketio.Conn, msg string) string
	OnLogin(s socketio.Conn, msg string) string
	OnLogout(s socketio.Conn, msg string) string
}

type IUsecase interface {
	SendMessage(userId uint, receiverId uint, message string) (Message, error)
	ReadAll(userId uint, receiverId uint) error

	Login(userId uint, socketId string) error
	Logout(userId uint, socketId string) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)

	IsOnline(userId uint) bool

	GetUserID(socketId string) (uint, error)
}

type IRepo interface {
	GetUserID(socketId string) (uint, error)
	IsLoggedIn(userId uint) bool
	Online(userId uint, socketId string)
	Offline(userId uint, socketId string)

	CreateMessage(msg Message) (Message, error)
	ReadAll(userId uint, senderId uint) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)
	//GetOverview(userId uint) ([]uint, error)
}

package chat

import (
	"time"

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
}

type IUsecase interface {
	SendMessage(userId uint, receiverId uint, message string) (Message, error)
	ReadAll(userId uint, receiverId uint) error

	Login(userId uint) ([]Message, error)
	Logout(userId uint) error
}

type IRepo interface {
	IsLoggedIn(userId uint) bool
	Online(userId uint, socketId string)
	Offline(userId uint)

	CreateMessage(msg Message) (Message, error)
	ReadAll(userId uint, senderId uint) error

	GetMessage(userId uint, user2Id uint) ([]Message, error)
}

package reminder

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/auth"
)

type Reminder struct {
	common.UUIDBase `gorm:"embedded;embeddedPrefix:reminder_"`
	Last            time.Time `gorm:"type:DATE"`
	Next            time.Time `gorm:"type:DATE"`
	ReminderType    string    `gorm:"type:VARCHAR(16);index:,type:brin;"`
	OwnerID         uint
	Owner           auth.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Reminder) TableName() string {
	return "reminder"
}

type IController interface {
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type IUsecase interface {
	GetOrCreate(user auth.User) ([]ReminderResponse, error)
	Update(user auth.User, item UpdateRequest) error
}

type IRepo interface {
	GetOrCreate(userId uint) ([]Reminder, error)
	Update(reminder Reminder) error
}

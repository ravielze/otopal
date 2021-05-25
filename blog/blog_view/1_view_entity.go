package blog_view

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/blog"
)

type (
	View struct {
		common.BigIntIDBase `gorm:"embedded;embeddedPrefix:blog_view_"`
		BlogID              string
		Blog                blog.Blog `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		IP                  string    `gorm:"type:VARCHAR(255);not null"`
		AccessTime          time.Time
	}

	BlogView struct {
		Blog  blog.Blog `gorm:"embedded"`
		Count int64
	}
)

func (View) TableName() string {
	return "view"
}

type IController interface {
	Top(ctx *gin.Context)
}

type IUsecase interface {
	AddView(blogId string, clientIp string) error
	Top(top int) ([]BlogView, error)
}

type IRepo interface {
	Create(view View) error
	GetLast(blogId string, clientIp string) (View, error)
	Top(top int) ([]BlogView, error)
}

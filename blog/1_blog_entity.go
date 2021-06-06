package blog

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/filemanager"
)

type Blog struct {
	common.UUIDBase `gorm:"embedded;embeddedPrefix:blog_"`
	common.InfoBase `gorm:"embedded;"`
	Title           string `gorm:"type:VARCHAR(256);"`
	Content         string `gorm:"type:VARCHAR(2097152);"`
	AuthorID        uint
	Author          auth.User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Thumbnails      []filemanager.File `gorm:"many2many:blog_thumbnails;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

const BLOG_PER_PAGE uint = 10

func (Blog) TableName() string {
	return "blog"
}

func (Blog) ThumbnailsTableName() string {
	return "blog_thumbnails"
}

type IController interface {
	GetBlogs(ctx *gin.Context)
	GetUserBlogs(ctx *gin.Context)
	GetBlog(ctx *gin.Context)

	Create(ctx *gin.Context)
	Edit(ctx *gin.Context)
	AddThumbnail(ctx *gin.Context)
	RemoveThumbnail(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type IUsecase interface {
	Create(user auth.User, item BlogRequest) (Blog, error)
	Edit(user auth.User, title string, lastEdit string, item BlogRequest) (Blog, error)

	GetBlogs(page uint) ([]Blog, error)
	GetUserBlogs(user auth.User, page uint) ([]Blog, error)
	GetBlog(title string, lastEdit string) (Blog, error)

	AddThumbnail(user auth.User, blogId string, item common.FileAttachment) error
	RemoveThumbnail(user auth.User, blogId string, fileId string) error

	Delete(user auth.User, blogId string) error
}

type IRepo interface {
	Create(blog Blog) (Blog, error)
	Edit(blog Blog) (Blog, error)

	GetBlogs(page uint) ([]Blog, error)
	GetUserBlogs(userId uint, page uint) ([]Blog, error)
	GetBlog(title string, lastEdit time.Time) (Blog, error)

	AddThumbnail(blog Blog, fileId string) error
	RemoveThumbnail(blog Blog, fileId string) error

	Delete(blog Blog) error
}

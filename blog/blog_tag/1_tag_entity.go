package blog_tag

import (
	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog"
)

type Tag struct {
	common.BigIntIDBase
	Name         string      `gorm:"type:VARCHAR(128);"`
	RelatedBlogs []blog.Blog `gorm:"many2many:blogs_tags;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (Tag) TableName() string {
	return "tag"
}

func (Tag) RelatedBlogsTableName() string {
	return "blogs_tags"
}

type IController interface {
	EditBlogTags(ctx *gin.Context)
	FindBlogs(ctx *gin.Context)
}

type IUsecase interface {
	EditBlogTags(user auth.User, blogId string, tags []string) error
	FindBlogs(tags []string) ([]blog.Blog, error)
}

type IRepo interface {
	CreateOrGet(tag Tag) (Tag, error)

	AddTag(userId uint, blogId string, tag Tag) error
	ClearTags(userId uint, blogId string) error

	FindBlog(tagName string) ([]blog.Blog, error)
}

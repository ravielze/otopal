package blog_connector

import "github.com/ravielze/otopal/auth"

type BlogTagClearUsecase interface {
	ClearTags(user auth.User, blogId string) error
}

var BTCU BlogTagClearUsecase

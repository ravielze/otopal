package blog_connector

type BlogViewUsecase interface {
	AddView(blogId string, clientIp string) error
}

var BVU BlogViewUsecase

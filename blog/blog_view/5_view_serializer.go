package blog_view

import "github.com/ravielze/otopal/blog"

type (
	BlogViewResponse struct {
		BlogResp blog.BlogResponse `json:"blog"`
		Count    int64             `json:"view_count"`
	}
)

func (item BlogView) Convert() BlogViewResponse {
	return BlogViewResponse{
		BlogResp: item.Blog.Convert(),
		Count:    item.Count,
	}
}

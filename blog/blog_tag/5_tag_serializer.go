package blog_tag

type TagRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

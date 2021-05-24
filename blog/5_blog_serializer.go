package blog

import (
	"strings"

	"github.com/ravielze/oculi/common/radix36"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/filemanager"
)

type (
	BlogRequest struct {
		Title   string `json:"title" binding:"required,lte=256,ascii"`
		Content string `json:"content" binding:"required,lte=2097152"`
	}

	BlogResponse struct {
		ID         string                     `json:"blog_id"`
		Title      string                     `json:"title"`
		Slug       string                     `json:"slug"`
		Content    string                     `json:"content"`
		Author     auth.UserResponse          `json:"author"`
		LastEdit   string                     `json:"last_edit"`
		Thumbnails []filemanager.FileResponse `json:"thumbnails"`
	}
)

func (item BlogRequest) Convert(authorId uint) Blog {
	return Blog{
		Title:    item.Title,
		Content:  item.Content,
		AuthorID: authorId,
	}
}

func (item Blog) Convert() BlogResponse {
	slug := strings.ReplaceAll(strings.ToLower(item.Title), " ", "-")
	thumbnails := make([]filemanager.FileResponse, len(item.Thumbnails))
	for i, x := range item.Thumbnails {
		thumbnails[i] = x.Convert()
	}
	if len(thumbnails) == 0 {
		thumbnails = append(thumbnails, filemanager.FileResponse{
			ID:       "default",
			FileType: "image/png",
		})
	}
	return BlogResponse{
		ID:         radix36.MustEncodeUUID(item.ID),
		Title:      item.Title,
		Slug:       slug,
		Content:    item.Content,
		Author:     item.Author.Convert(),
		LastEdit:   item.UpdatedAt.Format("02-01-2006"),
		Thumbnails: thumbnails,
	}
}

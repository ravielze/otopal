package blog

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	uc IUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc: uc,
	}
	blogGroup := g.Group("/blog")
	{
		blogGroup.GET("/", func(ctx *gin.Context) {
			fmt.Println("Module blog.")
		})
	}
	return cont
}

func (cont Controller) AddThumbnail(ctx *gin.Context){
    panic("not implemented")
}

func (cont Controller) Create(ctx *gin.Context){
    panic("not implemented")
}

func (cont Controller) Delete(ctx *gin.Context){
    panic("not implemented")
}

func (cont Controller) GetBlog(ctx *gin.Context){
    panic("not implemented")
}

func (cont Controller) GetBlogs(ctx *gin.Context){
    panic("not implemented")
}

func (cont Controller) GetUserBlogs(ctx *gin.Context){
    panic("not implemented")
}

func (cont Controller) RemoveThumbnail(ctx *gin.Context){
    panic("not implemented")
}
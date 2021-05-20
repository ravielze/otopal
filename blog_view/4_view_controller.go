package blog_view

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
	blog_viewGroup := g.Group("/blog_view")
	{
		blog_viewGroup.GET("/", func(ctx *gin.Context) {
			fmt.Println("Module blog_view.")
		})
	}
	return cont
}
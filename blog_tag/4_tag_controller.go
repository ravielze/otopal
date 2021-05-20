package blog_tag

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
	blog_tagGroup := g.Group("/blog_tag")
	{
		blog_tagGroup.GET("/", func(ctx *gin.Context) {
			fmt.Println("Module blog_tag.")
		})
	}
	return cont
}
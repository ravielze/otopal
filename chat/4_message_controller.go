package chat

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
	chatGroup := g.Group("/chat")
	{
		chatGroup.GET("/", func(ctx *gin.Context) {
			fmt.Println("Module chat.")
		})
	}
	return cont
}
package blog_view

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common/controller_utils"
	"github.com/ravielze/oculi/common/utils"
)

type Controller struct {
	uc IUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc: uc,
	}
	blogViewGroup := g.Group("/blog")
	{
		blogViewGroup.GET("/top", cont.Top)
	}
	return cont
}

func (cont Controller) Top(ctx *gin.Context) {
	ok, _, queries := controller_utils.NewControlChain(ctx).Query("top", "5").End()
	if ok {
		top, err := strconv.Atoi(queries["top"])
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}

		rawResult, err2 := cont.uc.Top(top)
		if err2 != nil {
			utils.AbortUsecaseError(ctx, err2)
			return
		}

		result := make([]BlogViewResponse, len(rawResult))
		for i, x := range rawResult {
			result[i] = x.Convert()
		}
		utils.OKAndResponseData(ctx, result)
	}
}

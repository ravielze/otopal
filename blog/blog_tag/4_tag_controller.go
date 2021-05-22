package blog_tag

import (
	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common/controller_utils"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/oculi/common/utils"
	"github.com/ravielze/otopal/auth"
)

type Controller struct {
	uc  IUsecase
	auc auth.IUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc:  uc,
		auc: module_manager.GetModule("auth").(auth.Module).Usecase(),
	}
	blog_tagGroup := g.Group("/blog/tags")
	blog_tagGroup.Use(cont.auc.AuthenticationNeeded(), cont.auc.AllowedRole(auth.ROLE_ADMIN))
	{
		blog_tagGroup.POST("/:blogid", cont.EditBlogTags)
	}
	return cont
}

func (cont Controller) EditBlogTags(ctx *gin.Context) {
	var obj TagRequest
	ok, params, _ := controller_utils.NewControlChain(ctx).ParamBase36ToUUID("blogid").BindJSON(&obj).End()
	if ok {
		user := cont.auc.GetUser(ctx)
		err := cont.uc.EditBlogTags(user, params["blogid"], obj.Tags)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponse(ctx)
	}
}

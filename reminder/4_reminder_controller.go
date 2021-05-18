package reminder

import (
	"github.com/gin-gonic/gin"
	cutils "github.com/ravielze/oculi/common/controller_utils"
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
	reminderGroup := g.Group("/reminder")
	reminderGroup.Use(cont.auc.AuthenticationNeeded(), cont.auc.AllowedRole(auth.ROLE_CUSTOMER))
	{
		reminderGroup.GET("/", cont.Get)
		reminderGroup.PUT("/", cont.Update)
	}
	return cont
}

func (cont Controller) Get(ctx *gin.Context) {
	user := cont.auc.GetUser(ctx)
	result, err := cont.uc.GetOrCreate(user)
	if err != nil {
		utils.AbortUsecaseError(ctx, err)
		return
	}
	utils.OKAndResponseData(ctx, result)
}

func (cont Controller) Update(ctx *gin.Context) {
	var obj UpdateRequest
	ok, _, _ := cutils.NewControlChain(ctx).BindJSON(&obj).End()
	if ok {
		user := cont.auc.GetUser(ctx)
		err := cont.uc.Update(user, obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponse(ctx)
	}
}

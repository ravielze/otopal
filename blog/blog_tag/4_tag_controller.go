package blog_tag

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common/controller_utils"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/oculi/common/utils"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog"
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

	publicBlogTagGroup := g.Group("/blog/tags")
	{
		publicBlogTagGroup.GET("/", cont.FindBlogs)
		publicBlogTagGroup.GET("/random", cont.RandomTags)
	}
	blogTagGroup := g.Group("/blog/tags")
	blogTagGroup.Use(cont.auc.AuthenticationNeeded(), cont.auc.AllowedRole(auth.ROLE_ADMIN))
	{
		blogTagGroup.POST("/:blogid", cont.EditBlogTags)
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

func (cont Controller) FindBlogs(ctx *gin.Context) {
	tags := ctx.QueryArray("tag")
	ok := true
	if len(tags) == 0 {
		ok = false
		utils.AbortUsecaseError(ctx, errors.New("parameter 'tag' is missing"))
		return
	}
	if ok {
		rawResult, err := cont.uc.FindBlogs(tags)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		result := make([]blog.BlogResponse, len(rawResult))
		for i, x := range rawResult {
			result[i] = x.Convert()
		}
		utils.OKAndResponseData(ctx, result)
	}
}

func (cont Controller) RandomTags(ctx *gin.Context) {
	ok, _, queries := controller_utils.NewControlChain(ctx).Query("amount", "6").End()
	if ok {
		amt, err := strconv.Atoi(queries["amount"])
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}

		rawResult, err2 := cont.uc.RandomTags(amt)
		if err2 != nil {
			utils.AbortUsecaseError(ctx, err2)
			return
		}

		result := make([]string, len(rawResult))
		for i, x := range rawResult {
			result[i] = x.Name
		}
		utils.OKAndResponseData(ctx, result)
	}
}

package blog

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	cutils "github.com/ravielze/oculi/common/controller_utils"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/oculi/common/utils"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog/blog_connector"
)

type Controller struct {
	uc   IUsecase
	auc  auth.IUsecase
	bvuc blog_connector.BlogViewUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc:   uc,
		auc:  module_manager.GetModule("auth").(auth.Module).Usecase(),
		bvuc: blog_connector.BVU,
	}
	blogGroup := g.Group("/blog")
	{
		blogGroup.GET("/blogs", cont.GetBlogs)
		blogGroup.GET("/info/:date/:slug", cont.GetBlog)
	}
	adminBlogGroup := g.Group("/blog")
	adminBlogGroup.Use(cont.auc.AuthenticationNeeded(), cont.auc.AllowedRole(auth.ROLE_ADMIN))
	{
		adminBlogGroup.POST("/create", cont.Create)
		adminBlogGroup.GET("/user", cont.GetUserBlogs)
		adminBlogGroup.DELETE("/delete/:blogid", cont.Delete)
		adminBlogGroup.POST("/thumbnail/:blogid", cont.AddThumbnail)
		adminBlogGroup.DELETE("/thumbnail/:blogid/:fileid", cont.RemoveThumbnail)
	}
	return cont
}

func (cont Controller) AddThumbnail(ctx *gin.Context) {
	var obj common.FileAttachment
	ok, params, _ := cutils.
		NewControlChain(ctx).
		BindForm(&obj).
		ParamBase36ToUUID("blogid").
		End()
	if ok {
		user := cont.auc.GetUser(ctx)
		err := cont.uc.AddThumbnail(user, params["blogid"], obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponse(ctx)
	}
}

func (cont Controller) Create(ctx *gin.Context) {
	var obj BlogRequest
	ok, _, _ := cutils.NewControlChain(ctx).BindJSON(&obj).End()
	if ok {
		user := cont.auc.GetUser(ctx)
		result, err := cont.uc.Create(user, obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponseData(ctx, result.Convert())
	}
}

func (cont Controller) Delete(ctx *gin.Context) {
	ok, params, _ := cutils.
		NewControlChain(ctx).
		ParamBase36ToUUID("blogid").
		End()
	if ok {
		user := cont.auc.GetUser(ctx)
		err := cont.uc.RemoveThumbnail(user, params["blogid"], params["fileid"])
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponse(ctx)
	}
}

func (cont Controller) GetBlog(ctx *gin.Context) {
	ok, params, _ := cutils.
		NewControlChain(ctx).
		Param("date").
		Param("slug").
		End()
	if ok {
		result, err := cont.uc.GetBlog(params["slug"], params["date"])
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		verr := cont.bvuc.AddView(result.ID, ctx.ClientIP())
		if verr != nil {
			utils.AbortUsecaseError(ctx, verr)
			return
		}
		utils.OKAndResponseData(ctx, result.Convert())
	}
}

func (cont Controller) GetBlogs(ctx *gin.Context) {
	ok, _, queries := cutils.
		NewControlChain(ctx).
		Query("page", "1").
		End()
	if ok {
		page, perr := strconv.Atoi(queries["page"])
		if perr != nil {
			utils.AbortUsecaseError(ctx, perr)
			return
		}

		rawResult, err := cont.uc.GetBlogs(uint(page))
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}

		result := make([]BlogResponse, len(rawResult))
		for i, x := range rawResult {
			result[i] = x.Convert()
		}
		utils.OKAndResponseData(ctx, result)
	}
}

func (cont Controller) GetUserBlogs(ctx *gin.Context) {
	ok, _, queries := cutils.
		NewControlChain(ctx).
		Query("page", "1").
		End()
	if ok {
		user := cont.auc.GetUser(ctx)
		page, perr := strconv.Atoi(queries["page"])
		if perr != nil {
			utils.AbortUsecaseError(ctx, perr)
			return
		}

		rawResult, err := cont.uc.GetUserBlogs(user, uint(page))
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}

		result := make([]BlogResponse, len(rawResult))
		for i, x := range rawResult {
			result[i] = x.Convert()
		}
		utils.OKAndResponseData(ctx, result)
	}
}

func (cont Controller) RemoveThumbnail(ctx *gin.Context) {
	ok, params, _ := cutils.
		NewControlChain(ctx).
		ParamBase36ToUUID("fileid").
		ParamBase36ToUUID("blogid").
		End()
	if ok {
		user := cont.auc.GetUser(ctx)
		err := cont.uc.RemoveThumbnail(user, params["blogid"], params["fileid"])
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponse(ctx)
	}
}

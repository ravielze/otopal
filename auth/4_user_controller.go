package auth

import (
	"github.com/gin-gonic/gin"
	cutils "github.com/ravielze/oculi/common/controller_utils"
	"github.com/ravielze/oculi/common/middleware"
	"github.com/ravielze/oculi/common/utils"
	"github.com/ravielze/otopal/chat/chat_connector"
)

type Controller struct {
	uc IUsecase
	cc chat_connector.ChatAuthUsecase
}

func NewController(g *gin.Engine, uc IUsecase) IController {
	cont := Controller{
		uc: uc,
		cc: chat_connector.CAU,
	}
	authGroup := g.Group("/auth")
	{
		authGroup.POST("/login", cont.Login)
		authGroup.POST("/register", cont.Register)
		authGroup.GET("/", uc.AuthenticationNeeded(), cont.Check)
		authGroup.PUT("/profile", uc.AuthenticationNeeded(), cont.Update)
		authGroup.POST("/registeradmin", middleware.GetStaticTokenMiddleware(), cont.RegisterAdmin)
	}
	g.GET("/technicians", uc.AuthenticationNeeded(), cont.GetTechnicians)
	return cont
}

func (cont Controller) Check(ctx *gin.Context) {
	utils.OKAndResponseData(ctx, cont.uc.GetUser(ctx).Convert())
}

func (cont Controller) Login(ctx *gin.Context) {
	var obj LoginRequest
	ok, _, _ := cutils.NewControlChain(ctx).BindJSON(&obj).End()
	if ok {
		result, err := cont.uc.Login(obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponseData(ctx, result)
	}
}

func (cont Controller) Register(ctx *gin.Context) {
	var obj RegisterRequest
	ok, _, _ := cutils.NewControlChain(ctx).BindJSON(&obj).End()
	if ok {
		result, err := cont.uc.Register(obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponseData(ctx, result)
	}
}

func (cont Controller) RegisterAdmin(ctx *gin.Context) {
	var obj RegisterRequest
	ok, _, _ := cutils.NewControlChain(ctx).BindJSON(&obj).End()
	if ok {
		result, err := cont.uc.RegisterAdmin(obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponseData(ctx, result)
	}
}

func (cont Controller) Update(ctx *gin.Context) {
	var obj UpdateRequest
	ok, _, _ := cutils.NewControlChain(ctx).BindJSON(&obj).End()
	if ok {
		user := cont.uc.GetUser(ctx)
		err := cont.uc.Update(user, obj)
		if err != nil {
			utils.AbortUsecaseError(ctx, err)
			return
		}
		utils.OKAndResponse(ctx)
	}
}

func (cont Controller) GetTechnicians(ctx *gin.Context) {
	technicians, err := cont.uc.GetTechnicians()
	if err != nil {
		utils.AbortUsecaseError(ctx, err)
		return
	}
	type Technician struct {
		User     UserResponse `json:"user"`
		IsOnline bool         `json:"online"`
	}
	result := make([]Technician, len(technicians))
	for i := range technicians {
		result[i] = Technician{
			User:     technicians[i].Convert(),
			IsOnline: cont.cc.IsOnline(technicians[i].ID),
		}
	}
	utils.OKAndResponseData(ctx, result)
}

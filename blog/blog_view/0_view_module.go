package blog_view

import (
	"github.com/gin-gonic/gin"
	"github.com/ravielze/otopal/blog/blog_connector"
	"gorm.io/gorm"
)

type Module struct {
	controller IController
	usecase    IUsecase
	repository IRepo
}

func (Module) Name() string {
	return "blog_view"
}

func (Module) Reset(db *gorm.DB) {
	db.Migrator().DropTable(&View{})
}

func (m Module) Usecase() IUsecase {
	return m.usecase
}

func NewModule(db *gorm.DB, g *gin.Engine) Module {
	repo := NewRepository(db)
	uc := NewUsecase(repo)
	cont := NewController(g, uc)

	db.AutoMigrate(&View{})
	blog_connector.BVU = uc

	return Module{
		controller: cont,
		usecase:    uc,
		repository: repo,
	}
}

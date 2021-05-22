package blog_view

import (
	"github.com/gin-gonic/gin"
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

func NewModule(db *gorm.DB, g *gin.Engine) Module {
	repo := NewRepository(db)
	uc := NewUsecase(repo)
	cont := NewController(g, uc)

	db.AutoMigrate(&View{})

	return Module{
		controller: cont,
		usecase:    uc,
		repository: repo,
	}
}

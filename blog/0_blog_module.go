package blog

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
	return "blog"
}

func (Module) Reset(db *gorm.DB) {
	db.Migrator().DropTable(&Blog{})
}

func (m Module) Usecase() IUsecase {
	return m.usecase
}

func NewModule(db *gorm.DB, g *gin.Engine) Module {
	repo := NewRepository(db)
	uc := NewUsecase(repo)
	cont := NewController(g, uc)

	db.AutoMigrate(&Blog{})

	return Module{
		controller: cont,
		usecase:    uc,
		repository: repo,
	}
}

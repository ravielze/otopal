package blog_tag

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
	return "blog_tag"
}

func (Module) Reset(db *gorm.DB) {
	db.Migrator().DropTable(&Tag{})
	db.Migrator().DropTable(Tag{}.RelatedBlogsTableName())
}

func NewModule(db *gorm.DB, g *gin.Engine) Module {
	repo := NewRepository(db)
	uc := NewUsecase(repo)
	cont := NewController(g, uc)

	db.AutoMigrate(&Tag{})

	return Module{
		controller: cont,
		usecase:    uc,
		repository: repo,
	}
}

package blog_tag

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
	blog_connector.BTCU = uc

	return Module{
		controller: cont,
		usecase:    uc,
		repository: repo,
	}
}

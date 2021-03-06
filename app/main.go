package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi"
	"github.com/ravielze/oculi/common/essentials"
	"github.com/ravielze/oculi/common/middleware"
	mm "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog"
	"github.com/ravielze/otopal/blog/blog_tag"
	"github.com/ravielze/otopal/blog/blog_view"
	"github.com/ravielze/otopal/chat"
	"github.com/ravielze/otopal/filemanager"
	"github.com/ravielze/otopal/reminder"
	"gorm.io/gorm"
)

func main() {
	oculi.New("Otopal", func(db *gorm.DB, g *gin.Engine) {
		middleware.InstallDefaultLimiter(g)
		// Add your middleware here
	}, func(db *gorm.DB, g *gin.Engine) {
		mm.AddModule(essentials.NewModule(db, g))
		mm.AddModule(auth.NewModule(db, g))
		mm.AddModule(filemanager.NewModule(db, g))
		mm.AddModule(reminder.NewModule(db, g))
		mm.AddModule(blog_view.NewModule(db, g))
		mm.AddModule(blog.NewModule(db, g))
		mm.AddModule(blog_tag.NewModule(db, g))
		mm.AddModule(chat.NewModule(db, g))
	}, func(db *gorm.DB, g *gin.Engine) {
		chatServer := chat.NewChatServer()
		chat.ChatServerInstance = chatServer
		middleware.InstallCors(g)
		signal.Notify(chatServer.Running, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		chatServer.Run(g)
		go g.Run()
		<-chatServer.Running
		os.Exit(0)
	})
	os.Exit(0)
}

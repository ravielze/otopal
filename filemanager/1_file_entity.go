package filemanager

import (
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/auth"
	"gorm.io/gorm"
)

type File struct {
	common.UUIDBase `gorm:"embedded;embeddedPrefix:file_"`
	common.InfoBase `gorm:"embedded"`
	FileGroup       string `gorm:"type:VARCHAR(256);index:,type:hash;"`
	FileType        string `gorm:"type:VARCHAR(256);"`
	FileExt         string `gorm:"type:VARCHAR(16);"`
	RealFilename    string `gorm:"type:VARCHAR(512);index:,sort:asc,type:btree"`
	Path            string `gorm:"type:VARCHAR(1024)"`
	Size            uint64 `gorm:"check:size >= 0;default:0"`
	OwnerID         uint
	Owner           auth.User `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (f *File) BeforeSave(db *gorm.DB) error {
	f.FileType = strings.ToLower(f.FileType)
	f.FileExt = strings.ToLower(f.FileExt)
	return nil
}

func (f *File) BeforeUpdate(db *gorm.DB) error {
	f.FileType = strings.ToLower(f.FileType)
	f.FileExt = strings.ToLower(f.FileExt)
	return nil
}

func (File) TableName() string {
	return "file"
}

type IController interface {
	GetFile(ctx *gin.Context)
	GetFilesByGroup(ctx *gin.Context)
}

type IUsecase interface {
	GetFile(idFile string) (FileResponse, error)
	GetRawFile(idFile string) (File, error)
	GetFilesByGroup(fileGroup string) ([]FileResponse, error)
	AddFile(user auth.User, fileGroup string, item common.FileAttachment) (FileResponse, error)
	DeleteFile(idFile string) error
}

type IRepo interface {
	GetFile(idFile string) (File, error)
	GetFilesByGroup(fileGroup string) ([]File, error)
	AddFile(userId uint, fileGroup string, attachment *multipart.FileHeader) (File, error)
	DeleteFile(idFile string) error
}

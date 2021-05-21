package blog

import (
	"strings"
	"time"

	"github.com/ravielze/oculi/common"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/filemanager"
)

type Usecase struct {
	repo IRepo
	fuc  filemanager.IUsecase
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{
		repo: repo,
		fuc:  module_manager.GetModule("filemanager").(filemanager.Module).Usecase(),
	}
}

func (uc Usecase) AddThumbnail(user auth.User, blogId string, item common.FileAttachment) error {
	fileResp, err := uc.fuc.AddFile(user, "blog", item)
	if err != nil {
		return err
	}

	err2 := uc.repo.AddThumbnail(Blog{
		UUIDBase: common.UUIDBase{
			ID: blogId,
		},
		AuthorID: user.ID,
	}, fileResp.ID)
	return err2
}

func (uc Usecase) RemoveThumbnail(user auth.User, blogId string, fileId string) error {
	err := uc.repo.RemoveThumbnail(Blog{
		UUIDBase: common.UUIDBase{
			ID: blogId,
		},
		AuthorID: user.ID,
	}, fileId)

	if err != nil {
		return err
	}

	err2 := uc.fuc.DeleteFile(fileId)
	return err2
}

func (uc Usecase) Create(user auth.User, item interface{}) (Blog, error) {
	panic("not implemented")
}

func (uc Usecase) Delete(user auth.User, blogId string) error {
	return uc.repo.Delete(Blog{
		UUIDBase: common.UUIDBase{
			ID: blogId,
		},
		AuthorID: user.ID,
	})
}

func (uc Usecase) GetBlog(title string, lastEdit string) (Blog, error) {
	timeParsed, err := time.Parse("02-01-2006", lastEdit)
	if err != nil {
		return Blog{}, err
	}
	titleTransformed := strings.ReplaceAll(title, "-", " ")
	return uc.repo.GetBlog(titleTransformed, timeParsed)
}

func (uc Usecase) GetBlogs(page uint) ([]Blog, error) {
	return uc.repo.GetBlogs(page)
}

func (uc Usecase) GetUserBlogs(user auth.User, page uint) ([]Blog, error) {
	return uc.repo.GetUserBlogs(user.ID, page)
}

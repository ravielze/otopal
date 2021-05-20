package blog

import (
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
	}, fileResp.ID)
	return err2
}

func (uc Usecase) Create(user auth.User, item interface{}) (Blog, error) {
	panic("not implemented")
}

func (uc Usecase) Delete(user auth.User, blogId string) error {
	panic("not implemented")
}

func (uc Usecase) GetBlog(title string, time string) (Blog, error) {
	panic("not implemented")
}

func (uc Usecase) GetBlogs(page uint) ([]Blog, error) {
	panic("not implemented")
}

func (uc Usecase) GetUserBlogs(user auth.User, page uint) ([]Blog, error) {
	panic("not implemented")
}

func (uc Usecase) RemoveThumbnail(user auth.User, blogId string, fileId string) error {
	panic("not implemented")
}

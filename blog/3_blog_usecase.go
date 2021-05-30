package blog

import (
	"errors"
	"strings"
	"time"

	"github.com/ravielze/oculi/common"
	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/oculi/common/radix36"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog/blog_connector"
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
	}, radix36.DecodeUUID(fileResp.ID).String())
	return err2
}

func (uc Usecase) RemoveThumbnail(user auth.User, blogId string, fileId string) error {
	err := uc.repo.RemoveThumbnail(Blog{
		UUIDBase: common.UUIDBase{
			ID: blogId,
		},
		AuthorID: user.ID,
	}, radix36.DecodeUUID(fileId).String())

	if err != nil {
		return err
	}

	err2 := uc.fuc.DeleteFile(fileId)
	return err2
}

func (uc Usecase) Delete(user auth.User, blogId string) error {
	//TODO remove file
	if err := blog_connector.BTCU.ClearTags(user, blogId); err != nil {
		return err
	}
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

func (uc Usecase) Create(user auth.User, item BlogRequest) (Blog, error) {
	if strings.Contains(item.Title, "-") {
		return Blog{}, errors.New("title cannot contains any '-'")
	}
	return uc.repo.Create(item.Convert(user.ID))
}

func (uc Usecase) Edit(user auth.User, title string, lastEdit string, item BlogRequest) (Blog, error) {
	oldBlog, err := uc.GetBlog(title, lastEdit)
	if err != nil {
		return Blog{}, err
	}
	blog := item.Convert(user.ID)
	blog.ID = oldBlog.ID
	return uc.repo.Edit(blog)
}

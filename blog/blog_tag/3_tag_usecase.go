package blog_tag

import (
	"errors"

	module_manager "github.com/ravielze/oculi/common/module"
	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog"
)

type Usecase struct {
	repo IRepo
	buc  blog.IUsecase
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{
		repo: repo,
		buc:  module_manager.GetModule("blog").(blog.Module).Usecase(),
	}
}

func (uc Usecase) EditBlogTags(user auth.User, blogId string, tags []string) error {
	tagsData := make([]Tag, len(tags))
	for i, x := range tags {
		if len(x) <= 0 || len(x) > 128 {
			continue
		}
		t, err := uc.repo.CreateOrGet(Tag{
			Name: x,
		})
		if err != nil {
			return err
		}
		tagsData[i] = t
	}

	err2 := uc.repo.ClearTags(user.ID, blogId)
	if err2 != nil {
		return err2
	}

	for _, x := range tagsData {
		err3 := uc.repo.AddTag(user.ID, blogId, x)
		if err3 != nil {
			return err3
		}
	}
	return nil
}

type B blog.Blog

func (uc Usecase) FindBlogs(tags []string) ([]blog.Blog, error) {
	tagsData := map[string][]blog.Blog{}
	var blogs []blog.Blog
	for _, x := range tags {
		if len(x) <= 0 || len(x) > 128 {
			continue
		}
		blog, err := uc.repo.FindBlog(x)
		if err != nil {
			return nil, err
		}
		if len(blog) > 0 {
			tagsData[x] = append(tagsData[x], blog...)
			blogs = append(blogs, blog...)
		}
	}

	blogCountMap := map[string]uint{}
	countBlogMap := map[uint][]blog.Blog{}
	blogMap := map[string]blog.Blog{}
	for _, blog := range blogs {
		blogCountMap[blog.ID]++
		blogMap[blog.ID] = blog
	}

	var maxVal uint = 0
	for k, v := range blogCountMap {
		if v > maxVal {
			maxVal = v
		}
		countBlogMap[v] = append(countBlogMap[v], blogMap[k])
	}

	var result []blog.Blog
	for i := uint(1); i <= maxVal; i++ {
		if len(countBlogMap[maxVal]) > 0 {
			result = append(result, countBlogMap[maxVal]...)
		}
	}
	return result, nil
}

func (uc Usecase) RandomTags(amount int) ([]Tag, error) {
	if amount > 30 || amount <= 0 {
		return nil, errors.New("amount is maxed at 30 and minimized at 0")
	}
	result, err := uc.repo.RandomTags(amount)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (uc Usecase) FindTag(title string, lastEdit string) ([]string, error) {
	blog, err := uc.buc.GetBlog(title, lastEdit)
	if err != nil {
		return nil, err
	}

	tags, err2 := uc.repo.FindTag(blog.ID)
	if err2 != nil {
		return nil, err2
	}
	result := make([]string, len(tags))
	for i, t := range tags {
		result[i] = t.Name
	}
	return result, nil
}

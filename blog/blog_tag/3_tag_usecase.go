package blog_tag

import (
	"fmt"

	"github.com/ravielze/otopal/auth"
	"github.com/ravielze/otopal/blog"
)

type Usecase struct {
	repo IRepo
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{repo: repo}
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
			fmt.Println(blog)
			tagsData[x] = append(tagsData[x], blog...)
			blogs = append(blogs, blog...)
		}
	}
	for i, j := range tagsData {
		fmt.Println(i, len(j))
	}
	fmt.Println(len(blogs))
	return nil, nil
}

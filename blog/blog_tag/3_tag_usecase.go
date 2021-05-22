package blog_tag

import "github.com/ravielze/otopal/auth"

type Usecase struct {
	repo IRepo
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{repo: repo}
}

func (uc Usecase) EditBlogTags(user auth.User, blogId string, tags []string) error {
	tagsData := make([]Tag, len(tags))
	for _, x := range tags {
		if len(x) <= 0 || len(x) > 128 {
			continue
		}

		tag, err := uc.repo.CreateOrGet(Tag{
			Name: x,
		})
		if err != nil {
			return err
		}
		tagsData = append(tagsData, tag)
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

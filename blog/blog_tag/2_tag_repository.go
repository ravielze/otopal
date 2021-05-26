package blog_tag

import (
	"fmt"
	"math/rand"

	"github.com/ravielze/otopal/blog"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepo {
	return Repository{db: db}
}

func (repo Repository) AddTag(userId uint, blogId string, tag Tag) error {
	var blog blog.Blog
	if err := repo.db.
		Where("author_id = ?", userId).
		Where("blog_id = ?", blogId).
		First(&blog).
		Error; err != nil {
		return err
	}

	if err2 := repo.db.
		Model(&tag).
		Omit("RelatedBlogs.*").
		Association("RelatedBlogs").
		Append(&blog); err2 != nil {
		return err2
	}
	return nil
}

func (repo Repository) ClearTags(userId uint, blogId string) error {
	var blog blog.Blog
	if err := repo.db.
		Where("author_id = ?", userId).
		Where("blog_id = ?", blogId).
		First(&blog).
		Error; err != nil {
		return err
	}

	if err2 := repo.db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE blog_id = '%s'", Tag{}.RelatedBlogsTableName(), blog.ID),
	).Error; err2 != nil {
		return err2
	}
	return nil
}

func (repo Repository) CreateOrGet(tag Tag) (Tag, error) {
	if err := repo.db.Where("name = ?", tag.Name).Attrs(Tag{Name: tag.Name}).FirstOrCreate(&tag).Error; err != nil {
		return Tag{}, err
	}
	return tag, nil
}

func (repo Repository) FindBlog(tagName string) ([]blog.Blog, error) {
	var result []blog.Blog
	var tagRelated Tag
	if err := repo.db.
		Where("name = ?", tagName).
		Preload("RelatedBlogs").
		Preload("RelatedBlogs.Author").
		Find(&tagRelated).
		Error; err != nil {
		return nil, err
	}
	result = append(result, tagRelated.RelatedBlogs...)
	return result, nil
}

func (repo Repository) RandomTags(amount int) ([]Tag, error) {
	var all []Tag
	if err := repo.db.Find(&all).Error; err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return all, nil
	}

	var result []Tag
	if amount >= len(all) {
		result = append(result, all...)
	} else {
		rand.Shuffle(len(all), func(i, j int) {
			all[i], all[j] = all[j], all[i]
		})
		result = append(result, all[:amount]...)
	}

	return result, nil
}

func (repo Repository) FindTag(blogId string) ([]Tag, error) {
	var tagIds []uint64
	if err := repo.db.
		Table(Tag{}.RelatedBlogsTableName()).
		Where("blog_id = ?", blogId).
		Find(&tagIds).
		Error; err != nil {
		return nil, err
	}
	var result []Tag
	if err := repo.db.
		Model(&Tag{}).
		Where("tag_id IN ?", tagIds).
		Order("name ASC").
		Find(&result).
		Error; err != nil {
		return nil, err
	}
	return result, nil
}

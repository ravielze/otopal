package blog_tag

import (
	"fmt"

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
		Association("Blog").
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
		fmt.Sprintf("DELETE FROM %s WHERE blog_id = %s", Tag{}.RelatedBlogsTableName(), blog.ID),
	).Error; err2 != nil {
		return err2
	}
	return nil
}

func (repo Repository) CreateOrGet(tag Tag) (Tag, error) {
	if err := repo.db.FirstOrCreate(&tag).Error; err != nil {
		return Tag{}, err
	}
	return tag, nil
}

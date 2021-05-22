package blog

import (
	"time"

	"github.com/ravielze/oculi/common"
	"github.com/ravielze/otopal/filemanager"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepo {
	return Repository{db: db}
}

func (repo Repository) AddThumbnail(blog Blog, fileId string) error {
	if err := repo.db.
		Model(&Blog{}).
		Where("blog_id = ?", blog.ID).
		Where("author_id = ?", blog.AuthorID).
		Omit("Thumbnails.*").
		Association("Thumbnails").
		Append(&filemanager.File{
			UUIDBase: common.UUIDBase{
				ID: fileId,
			},
		}); err != nil {
		return err
	}
	return nil
}

func (repo Repository) RemoveThumbnail(blog Blog, fileId string) error {
	if err := repo.db.
		Model(&Blog{}).
		Where("blog_id = ?", blog.ID).
		Where("author_id = ?", blog.AuthorID).
		Omit("Thumbnails.*").
		Association("Thumbnails").
		Delete(&filemanager.File{
			UUIDBase: common.UUIDBase{
				ID: fileId,
			},
		}); err != nil {
		return err
	}
	return nil
}

func (repo Repository) Create(blog Blog) (Blog, error) {
	if err := repo.db.
		Model(&Blog{}).
		Create(&blog).
		Error; err != nil {
		return Blog{}, err
	}
	repo.db.Model(&blog).Preload("Author").First(&blog)
	return blog, nil
}

func (repo Repository) Delete(blog Blog) error {
	if err := repo.db.
		Model(&Blog{}).
		Where("blog_id = ?", blog.ID).
		Where("author_id = ?", blog.AuthorID).
		Delete(&Blog{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo Repository) GetBlog(title string, lastEdit time.Time) (Blog, error) {
	var result Blog
	if err := repo.db.
		Preload("Author").
		Where("LOWER(title) LIKE LOWER(?)", title).
		Where("updated_at = ?", lastEdit).
		First(&result).Error; err != nil {
		return Blog{}, err
	}
	return result, nil
}

func (repo Repository) GetBlogs(page uint) ([]Blog, error) {
	var result []Blog
	pageOffset := (page - 1) * BLOG_PER_PAGE
	if err := repo.db.
		Preload("Author").
		Offset(int(pageOffset)).
		Limit(int(BLOG_PER_PAGE)).
		Find(&result).
		Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo Repository) GetUserBlogs(userId uint, page uint) ([]Blog, error) {
	var result []Blog
	pageOffset := (page - 1) * BLOG_PER_PAGE
	if err := repo.db.
		Preload("Author").
		Where("author_id = ?", userId).
		Offset(int(pageOffset)).
		Limit(int(BLOG_PER_PAGE)).
		Find(&result).
		Error; err != nil {
		return nil, err
	}
	return result, nil
}

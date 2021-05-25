package blog_view

import (
	"fmt"

	"github.com/ravielze/otopal/blog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepo {
	return Repository{db: db}
}

func (repo Repository) Create(view View) error {
	return repo.db.Create(&view).Error
}

func (repo Repository) GetLast(blogId string, clientIp string) (View, error) {
	var view View
	if err := repo.db.
		Order("access_time DESC").
		Where("blog_id = ?", blogId).
		Where("ip = ?", clientIp).
		Last(&view).
		Error; err != nil {
		return View{}, err
	}
	return view, nil
}

func (repo Repository) Top(top int) ([]BlogView, error) {
	var result []BlogView
	query := fmt.Sprintf(
		`SELECT T.id as blog_id,count 
FROM (SELECT blog_id as id, count(blog_id) AS count FROM "view" GROUP BY "blog_id" LIMIT %d) T
ORDER BY count desc`, top)
	if err := repo.db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}
	for i := range result {
		var blog blog.Blog
		if err := repo.db.
			Preload(clause.Associations).
			Find(&blog, "id = ?", result[i].Blog.ID).
			Error; err != nil {
			return nil, err
		}
		result[i].Blog = blog
	}
	return result, nil
}

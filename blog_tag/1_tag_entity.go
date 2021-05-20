package blog_tag

type Tag struct {
}

func (Tag) TableName() string {
	return "tag"
}

type IController interface {
}

type IUsecase interface {
}

type IRepo interface {
}
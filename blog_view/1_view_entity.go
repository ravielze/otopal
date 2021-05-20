package blog_view

type View struct {
}

func (View) TableName() string {
	return "view"
}

type IController interface {
}

type IUsecase interface {
}

type IRepo interface {
}
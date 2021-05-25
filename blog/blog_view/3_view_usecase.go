package blog_view

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Usecase struct {
	repo IRepo
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{repo: repo}
}

func (uc Usecase) AddView(blogId string, clientIp string) error {
	lastRecord, lastErr := uc.repo.GetLast(blogId, clientIp)
	if lastErr != nil {
		if !errors.Is(lastErr, gorm.ErrRecordNotFound) {
			return lastErr
		}
	} else {
		diff := time.Since(lastRecord.AccessTime).Minutes()
		if diff < 10 {
			return nil
		}
	}

	if err := uc.repo.Create(
		View{
			BlogID:     blogId,
			IP:         clientIp,
			AccessTime: time.Now(),
		},
	); err != nil {
		return err
	}
	return nil
}

func (uc Usecase) Top(top int) ([]BlogView, error) {
	if top > 20 || top <= 0 {
		return nil, errors.New("top is maxed at 20 and minimized at 0")
	}
	result, err := uc.repo.Top(top)
	if err != nil {
		return result, err
	}
	return result, nil
}

package reminder

import (
	"errors"
	"strings"

	"github.com/ravielze/otopal/auth"
)

type Usecase struct {
	repo IRepo
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{repo: repo}
}

func (uc Usecase) GetOrCreate(user auth.User) ([]ReminderResponse, error) {
	result, err := uc.repo.GetOrCreate(user.ID)
	if err != nil {
		return nil, err
	}

	resultConv := make([]ReminderResponse, len(result))
	for i, x := range result {
		resultConv[i] = x.Convert()
	}

	return resultConv, nil
}

func (uc Usecase) Update(user auth.User, item UpdateRequest) error {
	reminders, err := uc.repo.GetOrCreate(user.ID)
	if err != nil {
		return err
	}

	for _, x := range reminders {
		if strings.EqualFold(x.ReminderType, item.ReminderType) {
			reminder, err2 := item.Convert(x.ID, user.ID)
			if err2 != nil {
				return err2
			}
			uc.repo.Update(reminder)
			return nil
		}
	}
	return errors.New("update error")
}

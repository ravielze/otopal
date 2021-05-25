package chat

import (
	"errors"
	"time"
)

type Usecase struct {
	repo IRepo
}

func NewUsecase(repo IRepo) IUsecase {
	return Usecase{repo: repo}
}

func (uc Usecase) ReadAll(userId uint, receiverId uint) error {
	return uc.repo.ReadAll(userId, receiverId)
}

func (uc Usecase) SendMessage(userId uint, receiverId uint, message string) (Message, error) {
	if userId == 0 || receiverId == 0 {
		return Message{}, errors.New("senderId or receiverId cannot be 0")
	}
	return uc.repo.CreateMessage(Message{
		CreatedAt:  time.Now(),
		UserID:     userId,
		ReceiverID: receiverId,
		Message:    message,
		Read:       false,
	})
}

func (uc Usecase) Login(userId uint) ([]Message, error) {
	panic("not implemented")
}

func (uc Usecase) Logout(userId uint) error {
	panic("not implemented")
}

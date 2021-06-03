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
	if userId == 0 || receiverId == 0 {
		return errors.New("senderId or receiverId cannot be 0")
	}
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

func (uc Usecase) Login(userId uint, socketId int) error {
	uc.repo.Online(userId, socketId)
	return nil
}

func (uc Usecase) Logout(userId uint, socketId int) error {
	uc.repo.Offline(userId, socketId)
	return nil
}

func (uc Usecase) IsOnline(userId uint) bool {
	return uc.repo.IsLoggedIn(userId)
}

func (uc Usecase) GetMessage(userId uint, user2Id uint) ([]Message, error) {
	return uc.repo.GetMessage(userId, user2Id)
}

func (uc Usecase) GetUserID(socketId int) (uint, error) {
	return uc.repo.GetUserID(socketId)
}

func (uc Usecase) GetOverview(userId uint) ([]Message, []uint, error) {
	otherUsers, err := uc.repo.GetSender(userId)
	if err != nil {
		return nil, nil, err
	}
	msg := make([]Message, len(otherUsers))
	unread := make([]uint, len(otherUsers))
	for i, otherUserId := range otherUsers {
		ur, err2 := uc.repo.GetUnreadMessage(otherUserId, userId)
		if err2 != nil {
			return nil, nil, err2
		}
		m, err3 := uc.repo.GetLastMessage(userId, otherUserId)
		if err3 != nil {
			return nil, nil, err3
		}
		msg[i] = m
		unread[i] = ur
	}
	return msg, unread, nil
}

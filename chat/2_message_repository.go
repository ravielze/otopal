package chat

import (
	"sync"

	"gorm.io/gorm"
)

type OnlineData struct {
	sync.RWMutex
	// User ID -> Socket Connection ID
	userOnline *map[uint][]string
}

type Repository struct {
	db *gorm.DB
	od OnlineData
}

func NewRepository(db *gorm.DB) IRepo {
	return &Repository{
		db: db,
		od: OnlineData{
			userOnline: &map[uint][]string{},
		},
	}
}

func (repo *Repository) CreateMessage(msg Message) (Message, error) {
	if err := repo.db.
		Model(&Message{}).
		Create(&msg).
		Error; err != nil {
		return Message{}, err
	}
	repo.db.Model(&msg).Preload("User").Preload("Receiver").First(&msg)
	return msg, nil
}

func (repo *Repository) IsLoggedIn(userId uint) bool {
	repo.od.Lock()
	defer repo.od.Unlock()
	return len((*repo.od.userOnline)[userId]) != 0
}

func (repo *Repository) Offline(userId uint, socketId string) {
	repo.od.Lock()
	defer repo.od.Unlock()
	length := len((*repo.od.userOnline)[userId])
	if length == 1 {
		delete(*repo.od.userOnline, userId)
	} else {
		newArr := make([]string, length-1)
		for _, i := range (*repo.od.userOnline)[userId] {
			if i != socketId {
				newArr = append(newArr, i)
			}
		}
		(*repo.od.userOnline)[userId] = newArr
	}
}

func (repo *Repository) Online(userId uint, socketId string) {
	repo.od.Lock()
	defer repo.od.Unlock()
	(*repo.od.userOnline)[userId] = append((*repo.od.userOnline)[userId], socketId)
}

func (repo *Repository) ReadAll(userId uint, senderId uint) error {
	return repo.db.
		Model(&Message{}).
		Where("receiver_id = ?", userId).
		Where("user_id = ?", senderId).
		Update("read", true).
		Error
}

func (repo *Repository) GetMessage(userId uint, user2Id uint) ([]Message, error) {
	var result []Message
	if err := repo.db.Model(&Message{}).
		Where(map[string]interface{}{
			"user_id":     userId,
			"receiver_id": user2Id,
		}).
		Or(map[string]interface{}{
			"user_id":     user2Id,
			"receiver_id": userId,
		}).
		Order("created_time DESC").
		Find(&result).
		Error; err != nil {
		return nil, err
	}
	return result, nil
}

package chat_connector

type ChatAuthUsecase interface {
	IsOnline(userId uint) bool
}

var CAU ChatAuthUsecase

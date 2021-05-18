package auth

import "github.com/ravielze/oculi/common"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,lte=64"`
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email,ascii,lte=320"`
	Name        string `json:"name" binding:"required,ascii,lte=512"`
	PhoneNumber string `json:"phone_number" binding:"required,ascii,lte=24"`
	Password    string `json:"password" binding:"required,ascii,lte=64,gte=8"`
	Role        string `json:"role" binding:"required,lte=16"`
}

type UpdateRequest struct {
	Name        string `json:"name" binding:"required,ascii,lte=512"`
	PhoneNumber string `json:"phone_number" binding:"required,ascii,lte=24"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}

type UserTokenResponse struct {
	UserResponse
	Token string `json:"token"`
}

func (u User) ConvertToken(token string) UserTokenResponse {
	return UserTokenResponse{
		UserResponse: u.Convert(),
		Token:        token,
	}
}

func (u User) Convert() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Role:        Role(u.Role).String(),
	}
}

func (item RegisterRequest) Convert() User {
	return User{
		IntIDBase:   common.IntIDBase{},
		Email:       item.Email,
		Name:        item.Name,
		Password:    item.Password,
		PhoneNumber: item.PhoneNumber,
		Role:        int16(GetRole(item.Role)),
	}
}

func (item UpdateRequest) Convert(origin User) User {
	return User{
		//Unchangeable
		IntIDBase: common.IntIDBase{ID: origin.ID},
		Email:     origin.Email,
		Password:  origin.Password,
		Role:      origin.Role,

		//Changeable
		Name:        item.Name,
		PhoneNumber: item.PhoneNumber,
	}
}

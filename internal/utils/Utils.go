package utils

import (
	model2 "github.com/k0ch3gar/ozon-task/internal/graph/model"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

func FromDbUser(user *model.User) *model2.User {
	return &model2.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Deleted:   user.DeletedAt != nil,
	}
}

func FromApiUser(user *model2.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func FromUserInput(userInput *model2.UserInput) *model.User {
	return &model.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
}

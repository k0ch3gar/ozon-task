package service

import (
	"context"

	"github.com/k0ch3gar/ozon-task/internal/graph/model"
	"github.com/k0ch3gar/ozon-task/internal/storage"
	"github.com/k0ch3gar/ozon-task/internal/utils"
)

type UserService struct {
	us storage.UserStorage
}

func NewUserService(us storage.UserStorage) *UserService {
	return &UserService{
		us: us,
	}
}

func (us *UserService) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	user, err := us.us.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return utils.FromDbUser(user), err
}

func (us *UserService) CreateUser(ctx context.Context, userInput model.UserInput) (*model.User, error) {
	user := utils.FromUserInput(&userInput)
	err := us.us.InsertUser(user)
	if err != nil {
		return nil, err
	}

	return utils.FromDbUser(user), err
}

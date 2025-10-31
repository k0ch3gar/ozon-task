package utils

import (
	model2 "github.com/k0ch3gar/ozon-task/internal/graph/model"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

func FromStorageUser(user *model.User) *model2.User {
	if user.DeletedAt == nil {
		return &model2.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			Deleted:   false,
		}
	}

	dead := model2.DeadUser
	dead.ID = user.ID
	dead.CreatedAt = user.CreatedAt
	dead.Deleted = true

	return dead
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

func FromDbPost(post *model.Post) *model2.Post {
	if post.DeletedAt == nil {
		return &model2.Post{
			ID:            post.ID,
			Title:         post.Title,
			Body:          post.Body,
			AuthorID:      post.AuthorID,
			AllowComments: post.AllowComments,
			CreatedAt:     post.CreatedAt,
		}
	}

	return model2.DeadPost
}

func FromApiPost(post *model2.Post) *model.Post {
	return &model.Post{
		ID:            post.ID,
		Title:         post.Title,
		Body:          post.Body,
		AuthorID:      post.AuthorID,
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt,
	}
}

func FromPostInput(postInput *model2.PostInput) *model.Post {
	return &model.Post{
		Title:    postInput.Title,
		Body:     postInput.Body,
		AuthorID: &postInput.AuthorID,
	}
}

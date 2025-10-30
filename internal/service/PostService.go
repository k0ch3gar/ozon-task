package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/k0ch3gar/ozon-task/internal/config"
	"github.com/k0ch3gar/ozon-task/internal/graph/model"
	"github.com/k0ch3gar/ozon-task/internal/storage"
	"github.com/k0ch3gar/ozon-task/internal/utils"
)

type PostService struct {
	p        storage.PostStorage
	u        storage.UserStorage
	pageSize uint64
}

func NewPostService(params config.ApplicationParameters, p storage.PostStorage, u storage.UserStorage) *PostService {
	return &PostService{
		p:        p,
		u:        u,
		pageSize: params.PageSize,
	}
}

func (ps *PostService) GetPostsByPage(page uint64, ctx context.Context) ([]*model.Post, error) {
	posts, err := ps.p.GetFirstPostsFrom(page*ps.pageSize, ps.pageSize, ctx)
	if err != nil {
		return nil, err
	}

	apiPosts := make([]*model.Post, len(posts))
	for i := range apiPosts {
		apiPosts[i] = utils.FromDbPost(posts[i])
	}

	return apiPosts, nil
}

func (ps *PostService) GetPostByid(postId string, ctx context.Context) (*model.Post, error) {
	post, err := ps.p.GetPostById(postId, ctx)
	if err != nil {
		return nil, err
	}

	return utils.FromDbPost(post), err
}

func (ps *PostService) CreatePost(postInput model.PostInput, ctx context.Context) (*model.Post, error) {
	if ok, err := ps.u.ContainsById(postInput.AuthorID, ctx); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New(fmt.Sprintf("author does not exists: %s", postInput.AuthorID))
	}

	post := utils.FromPostInput(&postInput)
	err := ps.p.InsertPost(post, ctx)
	if err != nil {
		return nil, err
	}

	return utils.FromDbPost(post), err
}

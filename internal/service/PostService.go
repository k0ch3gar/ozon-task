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

	if len(posts) == 0 {
		return nil, nil
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

func (ps *PostService) UpdatePostTitle(ctx context.Context, postID string, title string) (*model.Post, error) {
	post, err := ps.GetPostByid(postID, ctx)
	if err != nil {
		return nil, err
	}

	post.Title = title
	err = ps.p.UpdatePost(utils.FromApiPost(post), ctx)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *PostService) UpdatePostBody(ctx context.Context, postID string, body string) (*model.Post, error) {
	post, err := ps.GetPostByid(postID, ctx)
	if err != nil {
		return nil, err
	}

	post.Body = body
	err = ps.p.UpdatePost(utils.FromApiPost(post), ctx)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *PostService) UpdatePostCommentsAllowance(ctx context.Context, postID string, allow bool) (*model.Post, error) {
	post, err := ps.GetPostByid(postID, ctx)
	if err != nil {
		return nil, err
	}

	post.AllowComments = allow
	err = ps.p.UpdatePost(utils.FromApiPost(post), ctx)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *PostService) DeletePost(ctx context.Context, postID string) (*string, error) {
	err := ps.p.DeletePost(postID, ctx)
	if err != nil {
		return nil, err
	}

	return &postID, nil
}

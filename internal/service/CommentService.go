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

type CommentService struct {
	u        storage.UserStorage
	p        storage.PostStorage
	c        storage.CommentStorage
	pageSize uint64
}

func NewCommentService(
	u storage.UserStorage,
	p storage.PostStorage,
	c storage.CommentStorage,
	params config.ApplicationParameters,
) *CommentService {
	return &CommentService{
		u:        u,
		c:        c,
		p:        p,
		pageSize: params.PageSize,
	}
}

func (cs *CommentService) GetPostCommentsByPage(postId string, page uint64, ctx context.Context) ([]*model.Comment, error) {
	post, err := cs.p.GetPostById(postId, ctx)
	if err != nil {
		return nil, errors.New("no such post")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed")
	}

	comments, err := cs.c.GetFirstCommentsByPost(postId, page*cs.pageSize, cs.pageSize, ctx)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, nil
	}

	apiComments := make([]*model.Comment, len(comments))
	for i := range apiComments {
		apiComments[i] = utils.FromStorageComment(comments[i])
	}

	return apiComments, nil
}

func (cs *CommentService) GetChildCommentsByPage(commentId string, page uint64, ctx context.Context) ([]*model.Comment, error) {
	comment, err := cs.GetCommentById(commentId, ctx)
	if err != nil {
		return nil, err
	}

	post, err := cs.p.GetPostById(comment.ParentPostID, ctx)
	if err != nil {
		return nil, errors.New("no such post")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed")
	}

	comments, err := cs.c.GetFirstCommentsByComment(commentId, page*cs.pageSize, cs.pageSize, ctx)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, nil
	}

	apiComments := make([]*model.Comment, len(comments))
	for i := range apiComments {
		apiComments[i] = utils.FromStorageComment(comments[i])
	}

	return apiComments, nil
}

func (cs *CommentService) CreateComment(commentInput model.CommentInput, ctx context.Context) (*model.Comment, error) {
	post, err := cs.p.GetPostById(commentInput.ParentPostID, ctx)
	if err != nil {
		return nil, errors.New("no such post")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed")
	}

	comment := utils.FromCommentInput(&commentInput)
	if ok, err := cs.u.ContainsById(*comment.AuthorID, ctx); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New(fmt.Sprintf("author does not exists: %s", *comment.AuthorID))
	}

	err = cs.c.InsertComment(comment, ctx)
	if err != nil {
		return nil, err
	}

	return utils.FromStorageComment(comment), nil
}

func (cs *CommentService) GetCommentById(commentId string, ctx context.Context) (*model.Comment, error) {
	comment, err := cs.c.GetCommentById(commentId, ctx)
	if err != nil {
		return nil, err
	}

	post, err := cs.p.GetPostById(comment.ParentPostID, ctx)
	if err != nil {
		return nil, errors.New("no such post")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed")
	}

	return utils.FromStorageComment(comment), nil
}

func (cs *CommentService) UpdateCommentBody(commentId string, body string, ctx context.Context) (*model.Comment, error) {
	comment, err := cs.c.GetCommentById(commentId, ctx)
	if err != nil {
		return nil, err
	}

	post, err := cs.p.GetPostById(comment.ParentPostID, ctx)
	if err != nil {
		return nil, errors.New("no such post")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed")
	}

	comment.Body = body
	err = cs.c.UpdateComment(comment, ctx)
	if err != nil {
		return nil, err
	}

	return utils.FromStorageComment(comment), nil
}

func (cs *CommentService) DeleteComment(commentId string, ctx context.Context) (*string, error) {
	comment, err := cs.c.GetCommentById(commentId, ctx)
	if err != nil {
		return nil, err
	}

	post, err := cs.p.GetPostById(comment.ParentPostID, ctx)
	if err != nil {
		return nil, errors.New("no such post")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed")
	}

	err = cs.c.DeleteComment(commentId, ctx)
	return &commentId, err
}

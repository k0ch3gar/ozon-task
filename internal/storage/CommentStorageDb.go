package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

type CommentStorageDb struct {
	db *pg.DB
	mu sync.Mutex
}

func NewDbCommentStorage(db *pg.DB) CommentStorage {
	return &CommentStorageDb{
		db: db,
		mu: sync.Mutex{},
	}
}

func (c *CommentStorageDb) GetCommentById(commentId string, ctx context.Context) (*model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	comment := &model.Comment{
		ID: commentId,
	}
	if err := getDataById(c.db, comment, ctx); err != nil {
		return nil, err
	}

	if comment == nil {
		return nil, errors.New("no such comment")
	}

	return comment, nil
}

func (c *CommentStorageDb) GetFirstCommentsByPost(postId string, offset, count uint64, ctx context.Context) ([]*model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var comments []*model.Comment
	query, err := buildQuery(c.db, comments, ctx)
	if err != nil {
		return nil, err
	}

	err = query.Where("parent_post_id = ?", postId).Order("created_at").Limit(int(count)).Offset(int(offset)).Select()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentStorageDb) GetFirstCommentsByComment(commentId string, offset, count uint64, ctx context.Context) ([]*model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var comments []*model.Comment
	query, err := buildQuery(c.db, comments, ctx)
	if err != nil {
		return nil, err
	}

	err = query.Where("parent_comment_id = ?", commentId).Order("created_at").Limit(int(count)).Offset(int(offset)).Select()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentStorageDb) InsertComment(comment *model.Comment, ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return insertData(c.db, comment, ctx)
}

func (c *CommentStorageDb) UpdateComment(newComment *model.Comment, ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return updateData(c.db, newComment, ctx)
}

func (c *CommentStorageDb) DeleteComment(commentId string, ctx context.Context) error {
	comment, err := c.GetCommentById(commentId, ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return deleteData(c.db, comment, ctx)
}

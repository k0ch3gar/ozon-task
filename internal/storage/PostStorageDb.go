package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

type PostStorageDb struct {
	mu sync.Mutex
	db *pg.DB
}

func NewDbPostStorage(db *pg.DB) PostStorage {
	return &PostStorageDb{
		db: db,
		mu: sync.Mutex{},
	}
}

func (p *PostStorageDb) GetFirstPostsFrom(offset uint64, count uint64, ctx context.Context) ([]*model.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var posts []*model.Post
	query, err := buildQuery(p.db, posts, ctx)
	if err != nil {
		return nil, err
	}

	err = query.Order("created_at").Limit(int(count)).Offset(int(offset)).Select()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostStorageDb) GetPostById(postId string, ctx context.Context) (*model.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	post := &model.Post{
		ID: postId,
	}
	if err := getDataById(p.db, post, ctx); err != nil {
		return nil, err
	}

	if post == nil {
		return nil, errors.New("no such post")
	}

	return post, nil
}

func (p *PostStorageDb) InsertPost(post *model.Post, ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return insertData(p.db, post, ctx)
}

func (p *PostStorageDb) UpdatePost(newPost *model.Post, ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return updateData(p.db, newPost, ctx)
}

func (p *PostStorageDb) DeletePost(postId string, ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return deleteDataById(p.db, postId, ctx)
}

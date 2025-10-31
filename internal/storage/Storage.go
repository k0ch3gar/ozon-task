package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/k0ch3gar/ozon-task/internal/config"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
	"go.uber.org/fx"
)

type UserStorage interface {
	GetUserById(userId string, ctx context.Context) (*model.User, error)
	InsertUser(user *model.User, ctx context.Context) error
	UpdateUser(newUser *model.User, ctx context.Context) error
	DeleteUser(userId string, ctx context.Context) (*model.User, error)
	ContainsByUsername(username string, ctx context.Context) (bool, error)
	ContainsById(userId string, ctx context.Context) (bool, error)
	GetUserByName(username string, ctx context.Context) (*model.User, error)
}

type PostStorage interface {
	GetFirstPostsFrom(offset uint64, count uint64, ctx context.Context) ([]*model.Post, error)
	GetPostById(postId string, ctx context.Context) (*model.Post, error)
	InsertPost(post *model.Post, ctx context.Context) error
	UpdatePost(newPost *model.Post, ctx context.Context) error
	DeletePost(postId string, ctx context.Context) error
}

type CommentStorage interface {
	GetCommentById(commentId string, ctx context.Context) (*model.Comment, error)
	GetFirstCommentsByPost(postId string, offset, count uint64, ctx context.Context) ([]*model.Comment, error)
	GetFirstCommentsByComment(commentId string, offset, count uint64, ctx context.Context) ([]*model.Comment, error)
	InsertComment(comment *model.Comment, ctx context.Context) error
	UpdateComment(newComment *model.Comment, ctx context.Context) error
	DeleteComment(commentId string, ctx context.Context) error
}

type StorageInMemoryShard[T any] struct {
	mu   sync.Mutex
	data map[string]*T
}

func NewDbConnection(opt pg.Options) (*pg.DB, error) {
	db := pg.Connect(&opt)
	if db == nil {
		return nil, errors.New(fmt.Sprintf("unable to connect to %s", opt.ToURL()))
	}

	return db, nil
}

func NewDbOpt() pg.Options {
	return pg.Options{
		Addr:     os.Getenv("PG_ADDR"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Database: os.Getenv("PG_DB"),
	}
}

func NewStorageModule(params config.ApplicationParameters) fx.Option {
	if params.PersistentStorageType {
		return fx.Module(
			`storage`,
			fx.Provide(
				NewDbOpt,
				NewDbConnection,
				NewDbUserStorage,
				NewDbPostStorage,
				NewDbCommentStorage,
			),
		)
	} else {
		return fx.Module(
			`storage`,
			fx.Provide(
				NewInMemoryUserStorage,
				NewInMemoryPostStorage,
				NewInMemoryCommentStorage,
			),
		)
	}
}

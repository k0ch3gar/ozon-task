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
	DeleteUser(userId string, ctx context.Context) error
}

type StorageDbShard struct {
	mu sync.Mutex
	db *pg.DB
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
		Database: os.Getenv("PG_NAME"),
	}
}

func NewStorageModule(params config.ApplicationParameters) fx.Option {
	if params.StorageType {
		return fx.Module(
			`storage`,
			fx.Provide(
				NewDbOpt,
				NewDbConnection,
				NewDbUserStorage,
			),
		)
	} else {
		return fx.Module(
			`storage`,
			fx.Provide(
				NewInMemoryUserStorage,
			),
		)
	}
}

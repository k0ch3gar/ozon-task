package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

const (
	ShardCount = 10
)

type UserStorageDb struct {
	shards         []*StorageDbShard
	insertionShard StorageDbShard
	shardCount     uint64
}

func NewDbUserStorage(db *pg.DB) UserStorage {
	shards := make([]*StorageDbShard, ShardCount)
	for i := range shards {
		shards[i] = &StorageDbShard{
			db: db,
			mu: sync.Mutex{},
		}
	}

	return &UserStorageDb{
		shardCount:     ShardCount,
		shards:         shards,
		insertionShard: StorageDbShard{db: db, mu: sync.Mutex{}},
	}
}

func (us *UserStorageDb) GetUserById(userId string, ctx context.Context) (*model.User, error) {
	uss, err := getStorageShard(us.shards, us.shardCount, userId)
	if err != nil {
		return nil, err
	}

	uss.mu.Lock()
	defer uss.mu.Unlock()

	user := &model.User{
		ID: userId,
	}
	if err := getDataById(uss.db, user, ctx); err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("no such user")
	}

	return user, nil
}

func (us *UserStorageDb) InsertUser(user *model.User, ctx context.Context) error {
	us.insertionShard.mu.Lock()
	defer us.insertionShard.mu.Unlock()

	return insertData(us.insertionShard.db, user, ctx)
}

func (us *UserStorageDb) UpdateUser(newUser *model.User, ctx context.Context) error {
	us.insertionShard.mu.Lock()
	defer us.insertionShard.mu.Unlock()

	return updateData(us.insertionShard.db, newUser, ctx)
}

func (us *UserStorageDb) DeleteUser(userId string, ctx context.Context) error {
	us.insertionShard.mu.Lock()
	defer us.insertionShard.mu.Unlock()

	return deleteDataById(us.insertionShard.db, userId, ctx)
}

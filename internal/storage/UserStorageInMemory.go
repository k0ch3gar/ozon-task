package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

type UserStorageInMemory struct {
	shards     []*StorageInMemoryShard[model.User]
	shardCount uint64
}

func NewInMemoryUserStorage() UserStorage {
	shards := make([]*StorageInMemoryShard[model.User], ShardCount)
	for i, _ := range shards {
		shards[i] = &StorageInMemoryShard[model.User]{}
		shards[i].mu = sync.Mutex{}
		shards[i].data = make(map[string]*model.User)
	}

	return &UserStorageInMemory{
		shardCount: ShardCount,
		shards:     shards,
	}
}

func (us *UserStorageInMemory) GetUserById(username string, ctx context.Context) (*model.User, error) {
	uss, err := getStorageShard(us.shards, us.shardCount, username)
	if err != nil {
		return nil, err
	}

	uss.mu.Lock()
	defer uss.mu.Unlock()

	user, ok := uss.data[username]
	if !ok {
		return nil, errors.New("no such user")
	}

	return user, nil
}

func (us *UserStorageInMemory) InsertUser(user *model.User, ctx context.Context) error {
	uss, err := getStorageShard(us.shards, us.shardCount, user.Username)
	if err != nil {
		return err
	}

	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[user.Username]
	if ok {
		return errors.New("such user exists")
	}

	uss.data[user.Username] = user
	return nil
}

func (us *UserStorageInMemory) UpdateUser(newUser *model.User, ctx context.Context) error {
	uss, err := getStorageShard(us.shards, us.shardCount, newUser.Username)
	if err != nil {
		return err
	}

	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[newUser.Username]
	if !ok {
		return errors.New(fmt.Sprintf("no such user with username: %s", newUser.Username))
	}

	uss.data[newUser.Username] = newUser
	return nil
}

func (us *UserStorageInMemory) DeleteUser(username string, ctx context.Context) error {
	uss, err := getStorageShard(us.shards, us.shardCount, username)
	if err != nil {
		return err
	}

	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[username]
	if !ok {
		return errors.New(fmt.Sprintf("no such user with username: %s", username))
	}

	delete(uss.data, username)
	return nil
}

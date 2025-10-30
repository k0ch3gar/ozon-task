package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/k0ch3gar/ozon-task/internal/config"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

type UserStorageInMemory struct {
	shards        []*StorageInMemoryShard[model.User]
	usernameShard map[string]uint64
	shardCount    uint64
	lastId        uint64
}

func NewInMemoryUserStorage(params config.ApplicationParameters) UserStorage {
	shards := make([]*StorageInMemoryShard[model.User], params.StorageShardsCount)
	for i, _ := range shards {
		shards[i] = &StorageInMemoryShard[model.User]{}
		shards[i].mu = sync.Mutex{}
		shards[i].data = make(map[string]*model.User)
	}

	return &UserStorageInMemory{
		shardCount:    params.StorageShardsCount,
		usernameShard: make(map[string]uint64),
		shards:        shards,
		lastId:        0,
	}
}

func (us *UserStorageInMemory) GetUserById(userId string, ctx context.Context) (*model.User, error) {
	idx, err := getStorageShardIdx(us.shards, us.shardCount, userId)
	if err != nil {
		return nil, err
	}

	uss := us.shards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	user, ok := uss.data[userId]
	if !ok {
		return nil, errors.New("no such user")
	}

	return user, nil
}

func (us *UserStorageInMemory) ContainsById(userId string, ctx context.Context) (bool, error) {
	idx, err := getStorageShardIdx(us.shards, us.shardCount, userId)
	if err != nil {
		return false, err
	}

	uss := us.shards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[userId]
	return ok, nil
}

func (us *UserStorageInMemory) ContainsByUsername(username string, ctx context.Context) (bool, error) {
	_, ok := us.usernameShard[username]
	return ok, nil
}

func (us *UserStorageInMemory) InsertUser(user *model.User, ctx context.Context) error {
	if ok, err := us.ContainsByUsername(user.Username, ctx); err != nil {
		return err
	} else if ok {
		return errors.New("user already exists")
	}

	id := strconv.FormatUint(us.lastId, 10)
	idx, err := getStorageShardIdx(us.shards, us.shardCount, id)
	if err != nil {
		return err
	}

	uss := us.shards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[user.ID]
	if ok {
		return errors.New("such user exists")
	}

	user.CreatedAt = time.Now().String()
	user.ID = id

	uss.data[user.ID] = user
	us.usernameShard[user.Username] = idx
	us.lastId++
	return nil
}

func (us *UserStorageInMemory) UpdateUser(newUser *model.User, ctx context.Context) error {
	idx, err := getStorageShardIdx(us.shards, us.shardCount, newUser.ID)
	if err != nil {
		return err
	}

	uss := us.shards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[newUser.ID]
	if !ok {
		return errors.New(fmt.Sprintf("no such user with id: %s", newUser.ID))
	}

	uss.data[newUser.ID] = newUser
	return nil
}

func (us *UserStorageInMemory) DeleteUser(userId string, ctx context.Context) error {
	idx, err := getStorageShardIdx(us.shards, us.shardCount, userId)
	if err != nil {
		return err
	}

	uss := us.shards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[userId]
	if !ok {
		return errors.New(fmt.Sprintf("no such user with userId: %s", userId))
	}

	delete(uss.data, userId)
	return nil
}

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
	idShards      []*StorageInMemoryShard[model.User]
	usernameShard []*StorageInMemoryShard[model.User]
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

	usernameShards := make([]*StorageInMemoryShard[model.User], params.StorageShardsCount)
	for i, _ := range shards {
		usernameShards[i] = &StorageInMemoryShard[model.User]{}
		usernameShards[i].mu = sync.Mutex{}
		usernameShards[i].data = make(map[string]*model.User)
	}

	return &UserStorageInMemory{
		shardCount:    params.StorageShardsCount,
		usernameShard: usernameShards,
		idShards:      shards,
		lastId:        0,
	}
}

func (us *UserStorageInMemory) GetUserById(userId string, ctx context.Context) (*model.User, error) {
	idx, err := getStorageShardIdx(us.idShards, us.shardCount, userId)
	if err != nil {
		return nil, err
	}

	uss := us.idShards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	user, ok := uss.data[userId]
	if !ok {
		return nil, errors.New("no such user")
	}

	if err = us.ValidateUserExistence(uss.data[userId]); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStorageInMemory) GetUserByName(username string, ctx context.Context) (*model.User, error) {
	idx, err := getStorageShardIdx(us.usernameShard, us.shardCount, username)
	if err != nil {
		return nil, err
	}

	usn := us.usernameShard[idx]
	usn.mu.Lock()
	defer usn.mu.Unlock()

	user, ok := usn.data[username]
	if !ok {
		return nil, errors.New("no such user")
	}

	if err = us.ValidateUserExistence(usn.data[username]); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStorageInMemory) ContainsById(userId string, ctx context.Context) (bool, error) {
	idx, err := getStorageShardIdx(us.idShards, us.shardCount, userId)
	if err != nil {
		return false, err
	}

	uss := us.idShards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[userId]
	if err = us.ValidateUserExistence(uss.data[userId]); err != nil {
		return false, err
	}
	return ok, nil
}

func (us *UserStorageInMemory) ContainsByUsername(username string, ctx context.Context) (bool, error) {
	idx, err := getStorageShardIdx(us.usernameShard, us.shardCount, username)
	if err != nil {
		return false, err
	}

	uss := us.usernameShard[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[username]
	if !ok {
		return false, nil
	}

	if err = us.ValidateUserExistence(uss.data[username]); err != nil {
		return false, err
	}

	return ok, nil
}

func (us *UserStorageInMemory) InsertUser(user *model.User, ctx context.Context) error {
	if ok, err := us.ContainsByUsername(user.Username, ctx); err != nil {
		return err
	} else if ok {
		return errors.New("user already exists")
	}

	id := strconv.FormatUint(us.lastId, 10)
	idx, err := getStorageShardIdx(us.idShards, us.shardCount, id)
	if err != nil {
		return err
	}

	uss := us.idShards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[user.ID]
	if ok {
		return errors.New("such user exists")
	}

	user.CreatedAt = time.Now().Format(time.RFC3339)
	user.ID = id

	uss.data[user.ID] = user
	idx, err = getStorageShardIdx(us.usernameShard, us.shardCount, user.Username)
	if err != nil {
		return err
	}

	usn := us.usernameShard[idx]
	usn.mu.Lock()
	defer usn.mu.Unlock()

	usn.data[user.Username] = user
	us.lastId++
	return nil
}

func (us *UserStorageInMemory) UpdateUser(newUser *model.User, ctx context.Context) error {
	idx, err := getStorageShardIdx(us.idShards, us.shardCount, newUser.ID)
	if err != nil {
		return err
	}

	uss := us.idShards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[newUser.ID]
	if !ok {
		return errors.New(fmt.Sprintf("no such user with id: %s", newUser.ID))
	}

	if err = us.ValidateUserExistence(uss.data[newUser.ID]); err != nil {
		return err
	}

	uss.data[newUser.ID] = newUser
	return nil
}

func (us *UserStorageInMemory) ValidateUserExistence(user *model.User) error {
	if user.DeletedAt != nil {
		return errors.New(fmt.Sprintf("user with this id is deleted: %s", user.ID))
	}

	return nil
}

func (us *UserStorageInMemory) DeleteUser(userId string, ctx context.Context) (*model.User, error) {
	idx, err := getStorageShardIdx(us.idShards, us.shardCount, userId)
	if err != nil {
		return nil, err
	}

	uss := us.idShards[idx]
	uss.mu.Lock()
	defer uss.mu.Unlock()

	_, ok := uss.data[userId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such user with id: %s", userId))
	}

	if uss.data[userId].DeletedAt != nil {
		return nil, errors.New(fmt.Sprintf("user with this id is already deleted: %s", userId))
	}

	deletionTime := time.Now().Format(time.RFC3339)
	uss.data[userId].DeletedAt = &deletionTime

	idx, err = getStorageShardIdx(us.usernameShard, us.shardCount, uss.data[userId].Username)
	if err != nil {
		return nil, err
	}

	usn := us.usernameShard[idx]
	usn.mu.Lock()
	defer usn.mu.Unlock()

	_, ok = usn.data[uss.data[userId].Username]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such user with username: %s", uss.data[userId].Username))
	}

	usn.data[uss.data[userId].Username].DeletedAt = &deletionTime
	return uss.data[userId], nil
}

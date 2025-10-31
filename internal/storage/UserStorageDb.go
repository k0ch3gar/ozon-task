package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

type UserStorageDb struct {
	mu sync.Mutex
	db *pg.DB
}

func (u *UserStorageDb) GetUserByName(username string, ctx context.Context) (*model.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user := &model.User{
		Username: username,
	}
	if err := getDataByUniqueColumn(u.db, user, "username", username, ctx); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errors.New("no such user")
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (u *UserStorageDb) ContainsByUsername(username string, ctx context.Context) (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user := &model.User{
		Username: username,
	}
	if err := getDataByUniqueColumn(u.db, user, "username", user.Username, ctx); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (u *UserStorageDb) ContainsById(userId string, ctx context.Context) (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user := &model.User{
		ID: userId,
	}
	if err := getDataById(u.db, user, ctx); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func NewDbUserStorage(db *pg.DB) UserStorage {
	return &UserStorageDb{
		db: db,
		mu: sync.Mutex{},
	}
}

func (u *UserStorageDb) GetUserById(userId string, ctx context.Context) (*model.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user := &model.User{
		ID: userId,
	}
	if err := getDataById(u.db, user, ctx); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errors.New("no such user")
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (u *UserStorageDb) InsertUser(user *model.User, ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return insertData(u.db, user, ctx)
}

func (u *UserStorageDb) UpdateUser(newUser *model.User, ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return updateData(u.db, newUser, ctx)
}

func (u *UserStorageDb) DeleteUser(userId string, ctx context.Context) (*model.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user := &model.User{
		ID: userId,
	}
	if err := getDataById(u.db, user, ctx); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errors.New("no such user")
		} else {
			return nil, err
		}
	}

	return user, deleteData(u.db, user, ctx)
}

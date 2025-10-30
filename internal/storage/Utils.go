package storage

import (
	"context"
	"errors"
	"hash/maphash"

	"github.com/go-pg/pg/v10"
)

var seed = maphash.MakeSeed()

func getShardIndex(key string, count uint64) uint64 {
	return maphash.String(seed, key) % count
}

func getStorageShardIdx[T any](shards []*T, shardCount uint64, id string) (uint64, error) {
	idx := getShardIndex(id, shardCount)
	if idx < 0 || idx >= shardCount {
		return 0, errors.New("shard index out of range")
	}

	uss := shards[idx]
	if uss == nil {
		return 0, errors.New("shard is nil")
	}

	return idx, nil
}

func getDataByUniqueColumn(db *pg.DB, data interface{}, column string, value string, ctx context.Context) error {
	query, err := buildQuery(db, data, ctx)
	if err != nil {
		return err
	}

	return query.Where(column+" = ?", value).Select()
}

func getDataById(db *pg.DB, data interface{}, ctx context.Context) error {
	query, err := buildQuery(db, data, ctx)
	if err != nil {
		return err
	}

	return query.WherePK().Select()
}

func insertData(db *pg.DB, data interface{}, ctx context.Context) error {
	query, err := buildQuery(db, data, ctx)
	if err != nil {
		return err
	}

	_, err = query.Insert()
	return err
}

func updateData(db *pg.DB, data interface{}, ctx context.Context) error {
	query, err := buildQuery(db, data, ctx)
	if err != nil {
		return err
	}

	_, err = query.WherePK().Update()
	return err
}

func deleteDataById(db *pg.DB, data interface{}, ctx context.Context) error {
	query, err := buildQuery(db, data, ctx)
	if err != nil {
		return err
	}

	_, err = query.Delete()
	return err
}

func buildQuery(db *pg.DB, data interface{}, ctx context.Context) (*pg.Query, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return db.WithContext(ctx).Model(data), nil
}

package storage

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/k0ch3gar/ozon-task/internal/config"
	"github.com/k0ch3gar/ozon-task/internal/storage/model"
)

type CommentStorageInMemory struct {
	shards     []*StorageInMemoryShard[model.Comment]
	shardCount uint64
	lastId     uint64
}

func NewInMemoryCommentStorage(params config.ApplicationParameters) CommentStorage {
	shards := make([]*StorageInMemoryShard[model.Comment], params.StorageShardsCount)
	for i := range shards {
		shards[i] = &StorageInMemoryShard[model.Comment]{}
		shards[i].mu = sync.Mutex{}
		shards[i].data = make(map[string]*model.Comment)
	}

	return &CommentStorageInMemory{
		shardCount: params.StorageShardsCount,
		shards:     shards,
		lastId:     0,
	}
}

func (c *CommentStorageInMemory) GetCommentById(commentId string, ctx context.Context) (*model.Comment, error) {
	idx, err := getStorageShardIdx(c.shards, c.shardCount, commentId)
	if err != nil {
		return nil, err
	}

	cs := c.shards[idx]
	cs.mu.Lock()
	defer cs.mu.Unlock()

	comment, ok := cs.data[commentId]
	if !ok {
		return nil, errors.New("no such comment")
	}

	if err = c.ValidateCommentExistence(cs.data[commentId]); err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *CommentStorageInMemory) ValidateCommentExistence(comment *model.Comment) error {
	if comment.DeletedAt != nil {
		return errors.New(fmt.Sprintf("comment with this id is deleted: %s", comment.ID))
	}

	return nil
}

func (c *CommentStorageInMemory) GetFirstCommentsByPost(postId string, offset, count uint64, ctx context.Context) ([]*model.Comment, error) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	var allComments []*model.Comment

	for _, shard := range c.shards {
		wg.Add(1)
		go func(cs *StorageInMemoryShard[model.Comment]) {
			defer wg.Done()

			cs.mu.Lock()
			defer cs.mu.Unlock()

			var comments []*model.Comment
			for _, val := range cs.data {
				if val.ParentPostID != postId || val.ParentCommentID != nil {
					continue
				}

				comments = append(comments, val)
			}

			mu.Lock()
			allComments = append(allComments, comments...)
			mu.Unlock()
		}(shard)
	}

	wg.Wait()

	sort.Slice(allComments, func(i, j int) bool {

		t1, err := time.Parse(time.RFC3339, allComments[i].CreatedAt)
		if err != nil {
			panic(err)
		}

		t2, _ := time.Parse(time.RFC3339, allComments[j].CreatedAt)
		if err != nil {
			panic(err)
		}

		return t1.Before(t2)
	})

	if offset >= uint64(len(allComments)) {
		return []*model.Comment{}, nil
	}

	end := offset + count
	if end > uint64(len(allComments)) {
		end = uint64(len(allComments))
	}

	return allComments[offset:end], nil
}

func (c *CommentStorageInMemory) GetFirstCommentsByComment(commentId string, offset, count uint64, ctx context.Context) ([]*model.Comment, error) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	var allComments []*model.Comment

	for _, shard := range c.shards {
		wg.Add(1)
		go func(cs *StorageInMemoryShard[model.Comment]) {
			defer wg.Done()

			cs.mu.Lock()
			defer cs.mu.Unlock()

			var comments []*model.Comment
			for _, val := range cs.data {
				if val.ParentCommentID == nil || *val.ParentCommentID != commentId {
					continue
				}

				comments = append(comments, val)
			}

			mu.Lock()
			allComments = append(allComments, comments...)
			mu.Unlock()
		}(shard)
	}

	wg.Wait()

	sort.Slice(allComments, func(i, j int) bool {

		t1, err := time.Parse(time.RFC3339, allComments[i].CreatedAt)
		if err != nil {
			panic(err)
		}

		t2, _ := time.Parse(time.RFC3339, allComments[j].CreatedAt)
		if err != nil {
			panic(err)
		}

		return t1.Before(t2)
	})

	if offset >= uint64(len(allComments)) {
		return []*model.Comment{}, nil
	}

	end := offset + count
	if end > uint64(len(allComments)) {
		end = uint64(len(allComments))
	}

	return allComments[offset:end], nil
}

func (c *CommentStorageInMemory) InsertComment(comment *model.Comment, ctx context.Context) error {
	id := strconv.FormatUint(c.lastId, 10)
	idx, err := getStorageShardIdx(c.shards, c.shardCount, id)
	if err != nil {
		return err
	}

	cs := c.shards[idx]
	cs.mu.Lock()
	defer cs.mu.Unlock()

	_, ok := cs.data[id]
	if ok {
		return errors.New("such comment already exists")
	}

	comment.ID = id
	comment.CreatedAt = time.Now().Format(time.RFC3339)

	cs.data[id] = comment
	c.lastId++
	return nil
}

func (c *CommentStorageInMemory) UpdateComment(newComment *model.Comment, ctx context.Context) error {
	idx, err := getStorageShardIdx(c.shards, c.shardCount, newComment.ID)
	if err != nil {
		return err
	}

	cs := c.shards[idx]
	cs.mu.Lock()
	defer cs.mu.Unlock()

	_, ok := cs.data[newComment.ID]
	if !ok {
		return errors.New("no such comment exists")
	}

	if err = c.ValidateCommentExistence(newComment); err != nil {
		return err
	}

	cs.data[newComment.ID] = newComment
	return nil
}

func (c *CommentStorageInMemory) DeleteComment(commentId string, ctx context.Context) error {
	idx, err := getStorageShardIdx(c.shards, c.shardCount, commentId)
	if err != nil {
		return err
	}

	cs := c.shards[idx]
	cs.mu.Lock()
	defer cs.mu.Unlock()

	comment, ok := cs.data[commentId]
	if !ok {
		return errors.New("no such comment exists")
	}

	if err = c.ValidateCommentExistence(comment); err != nil {
		return err
	}

	deletionTime := time.Now().Format(time.RFC3339)
	cs.data[commentId].DeletedAt = &deletionTime
	return nil
}

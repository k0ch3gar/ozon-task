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

type PostStorageInMemory struct {
	shards     []*StorageInMemoryShard[model.Post]
	shardCount uint64
	lastId     uint64
}

func NewInMemoryPostStorage(params config.ApplicationParameters) PostStorage {
	shards := make([]*StorageInMemoryShard[model.Post], params.StorageShardsCount)
	for i := range shards {
		shards[i] = &StorageInMemoryShard[model.Post]{}
		shards[i].mu = sync.Mutex{}
		shards[i].data = make(map[string]*model.Post)
	}

	return &PostStorageInMemory{
		shardCount: params.StorageShardsCount,
		shards:     shards,
		lastId:     0,
	}
}

func (p *PostStorageInMemory) GetFirstPostsFrom(offset uint64, count uint64, ctx context.Context) ([]*model.Post, error) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	var allPosts []*model.Post

	for _, shard := range p.shards {
		wg.Add(1)
		go func(ps *StorageInMemoryShard[model.Post]) {
			defer wg.Done()

			ps.mu.Lock()
			defer ps.mu.Unlock()

			var posts []*model.Post
			for _, val := range ps.data {
				posts = append(posts, val)
			}

			mu.Lock()
			allPosts = append(allPosts, posts...)
			mu.Unlock()
		}(shard)
	}

	wg.Wait()

	sort.Slice(allPosts, func(i, j int) bool {

		t1, err := time.Parse(time.RFC3339, allPosts[i].CreatedAt)
		if err != nil {
			panic(err)
		}

		t2, _ := time.Parse(time.RFC3339, allPosts[j].CreatedAt)
		if err != nil {
			panic(err)
		}

		return t1.Before(t2)
	})

	if offset >= uint64(len(allPosts)) {
		return []*model.Post{}, nil
	}

	end := offset + count
	if end > uint64(len(allPosts)) {
		end = uint64(len(allPosts))
	}

	return allPosts[offset:end], nil
}

func (p *PostStorageInMemory) GetPostById(postId string, ctx context.Context) (*model.Post, error) {
	idx, err := getStorageShardIdx(p.shards, p.shardCount, postId)
	if err != nil {
		return nil, err
	}

	ps := p.shards[idx]
	ps.mu.Lock()
	defer ps.mu.Unlock()

	post, ok := ps.data[postId]
	if !ok {
		return nil, errors.New("no such post")
	}

	if err = p.ValidatePostExistence(ps.data[postId]); err != nil {
		return nil, err
	}

	return post, nil
}

func (p *PostStorageInMemory) InsertPost(post *model.Post, ctx context.Context) error {
	id := strconv.FormatUint(p.lastId, 10)
	idx, err := getStorageShardIdx(p.shards, p.shardCount, id)
	if err != nil {
		return err
	}

	ps := p.shards[idx]
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.data[id]
	if ok {
		return errors.New("such post already exists")
	}

	post.ID = id
	post.CreatedAt = time.Now().Format(time.RFC3339)

	ps.data[id] = post
	p.lastId++
	return nil
}

func (p *PostStorageInMemory) UpdatePost(newPost *model.Post, ctx context.Context) error {
	idx, err := getStorageShardIdx(p.shards, p.shardCount, newPost.ID)
	if err != nil {
		return err
	}

	ps := p.shards[idx]
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.data[newPost.ID]
	if !ok {
		return errors.New(fmt.Sprintf("no such post with id: %s", newPost.ID))
	}

	if err = p.ValidatePostExistence(ps.data[newPost.ID]); err != nil {
		return err
	}

	ps.data[newPost.ID] = newPost
	return nil
}

func (p *PostStorageInMemory) ValidatePostExistence(post *model.Post) error {
	if post.DeletedAt != nil {
		return errors.New(fmt.Sprintf("post with this id is deleted: %s", post.ID))
	}

	return nil
}

func (p *PostStorageInMemory) DeletePost(postId string, ctx context.Context) error {
	idx, err := getStorageShardIdx(p.shards, p.shardCount, postId)
	if err != nil {
		return err
	}

	ps := p.shards[idx]
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.data[postId]
	if !ok {
		return errors.New(fmt.Sprintf("no such post with id: %s", postId))
	}

	if ps.data[postId].DeletedAt != nil {
		return errors.New(fmt.Sprintf("post with this id is already deleted: %s", postId))
	}

	deletionTime := time.Now().Format(time.RFC3339)
	ps.data[postId].DeletedAt = &deletionTime
	return nil
}

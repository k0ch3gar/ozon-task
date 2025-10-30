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
			posts := make([]*model.Post, len(ps.data))
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
		return allPosts[i].CreatedAt > allPosts[j].CreatedAt
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
	post.CreatedAt = time.Now().String()

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

	ps.data[newPost.ID] = newPost
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

	delete(ps.data, postId)
	return nil
}

package graph

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/k0ch3gar/ozon-task/internal/config"
	"github.com/k0ch3gar/ozon-task/internal/graph/model"
	"github.com/k0ch3gar/ozon-task/internal/service"
	"github.com/k0ch3gar/ozon-task/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUserCreated(t *testing.T) {
	params := config.ApplicationParameters{
		StorageShardsCount:    6,
		Port:                  "8080",
		PersistentStorageType: false,
		Debug:                 true,
		PageSize:              1,
	}

	u := storage.NewInMemoryUserStorage(params)
	p := storage.NewInMemoryPostStorage(params)
	c := storage.NewInMemoryCommentStorage(params)

	resolver := NewResolver(
		service.NewUserService(
			u,
		),
		service.NewPostService(
			params,
			p,
			u,
		),
		service.NewCommentService(
			u,
			c,
			params,
		),
		service.NewSubscriptionService(),
	)

	userInput := model.UserInput{
		Username: "foo",
		Email:    "bar",
		Password: "baz",
	}

	ctx := context.Background()
	user, err := resolver.Mutation().CreateUser(ctx, userInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, user.Username, userInput.Username)
	assert.Equal(t, user.Email, userInput.Email)
}

func TestUserExistence(t *testing.T) {
	params := config.ApplicationParameters{
		StorageShardsCount:    6,
		Port:                  "8080",
		PersistentStorageType: false,
		Debug:                 true,
		PageSize:              1,
	}

	u := storage.NewInMemoryUserStorage(params)
	p := storage.NewInMemoryPostStorage(params)
	c := storage.NewInMemoryCommentStorage(params)

	resolver := NewResolver(
		service.NewUserService(
			u,
		),
		service.NewPostService(
			params,
			p,
			u,
		),
		service.NewCommentService(
			u,
			c,
			params,
		),
		service.NewSubscriptionService(),
	)

	userInput := model.UserInput{
		Username: "foo",
		Email:    "bar",
		Password: "baz",
	}

	ctx := context.Background()
	user, err := resolver.Mutation().CreateUser(ctx, userInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	user, err = resolver.Query().UserByID(ctx, user.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, user.Username, userInput.Username)
	assert.Equal(t, user.Email, userInput.Email)
}

func TestPostCreationAndExistence(t *testing.T) {
	params := config.ApplicationParameters{
		StorageShardsCount:    6,
		Port:                  "8080",
		PersistentStorageType: false,
		Debug:                 true,
		PageSize:              1,
	}

	u := storage.NewInMemoryUserStorage(params)
	p := storage.NewInMemoryPostStorage(params)
	c := storage.NewInMemoryCommentStorage(params)

	resolver := NewResolver(
		service.NewUserService(
			u,
		),
		service.NewPostService(
			params,
			p,
			u,
		),
		service.NewCommentService(
			u,
			c,
			params,
		),
		service.NewSubscriptionService(),
	)

	userInput := model.UserInput{
		Username: "foo",
		Email:    "bar",
		Password: "baz",
	}

	ctx := context.Background()
	user, err := resolver.Mutation().CreateUser(ctx, userInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	postInput := model.PostInput{
		AuthorID: user.ID,
		Title:    "title1",
		Body:     "body1",
	}

	post, err := resolver.Mutation().CreatePost(ctx, postInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, *post.AuthorID, user.ID)
	assert.Equal(t, post.Body, postInput.Body)
	assert.Equal(t, post.Title, postInput.Title)

	post, err = resolver.Query().Post(ctx, post.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, *post.AuthorID, user.ID)
	assert.Equal(t, post.Body, postInput.Body)
	assert.Equal(t, post.Title, postInput.Title)
}

func TestCommentCreationAndExistence(t *testing.T) {
	params := config.ApplicationParameters{
		StorageShardsCount:    6,
		Port:                  "8080",
		PersistentStorageType: false,
		Debug:                 true,
		PageSize:              1,
	}

	u := storage.NewInMemoryUserStorage(params)
	p := storage.NewInMemoryPostStorage(params)
	c := storage.NewInMemoryCommentStorage(params)

	resolver := NewResolver(
		service.NewUserService(
			u,
		),
		service.NewPostService(
			params,
			p,
			u,
		),
		service.NewCommentService(
			u,
			c,
			params,
		),
		service.NewSubscriptionService(),
	)

	userInput := model.UserInput{
		Username: "foo",
		Email:    "bar",
		Password: "baz",
	}

	ctx := context.Background()
	user, err := resolver.Mutation().CreateUser(ctx, userInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	postInput := model.PostInput{
		AuthorID: user.ID,
		Title:    "title1",
		Body:     "body1",
	}

	post, err := resolver.Mutation().CreatePost(ctx, postInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	commentInput := model.CommentInput{
		AuthorID:     user.ID,
		ParentPostID: post.ID,
		Body:         "body2",
	}

	comment, err := resolver.Mutation().CreateComment(ctx, commentInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, *comment.AuthorID, user.ID)
	assert.Equal(t, comment.Body, commentInput.Body)
	assert.Equal(t, comment.ParentPostID, post.ID)
	assert.True(t, comment.ParentCommentID == nil)
}
func TestCommentSubscription(t *testing.T) {
	params := config.ApplicationParameters{
		StorageShardsCount:    6,
		Port:                  "8080",
		PersistentStorageType: false,
		Debug:                 true,
		PageSize:              1,
	}

	u := storage.NewInMemoryUserStorage(params)
	p := storage.NewInMemoryPostStorage(params)
	c := storage.NewInMemoryCommentStorage(params)

	resolver := NewResolver(
		service.NewUserService(
			u,
		),
		service.NewPostService(
			params,
			p,
			u,
		),
		service.NewCommentService(
			u,
			c,
			params,
		),
		service.NewSubscriptionService(),
	)

	userInput := model.UserInput{
		Username: "foo",
		Email:    "bar",
		Password: "baz",
	}

	ctx := context.Background()
	user, err := resolver.Mutation().CreateUser(ctx, userInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	postInput := model.PostInput{
		AuthorID: user.ID,
		Title:    "title1",
		Body:     "body1",
	}

	post, err := resolver.Mutation().CreatePost(ctx, postInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	commentInput := model.CommentInput{
		AuthorID:     user.ID,
		ParentPostID: post.ID,
		Body:         "body2",
	}

	ch, err := resolver.Subscription().CommentCreated(ctx, post.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(time.Second * 3)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case comment := <-ch:
				assert.Equal(t, *comment.AuthorID, user.ID)
				assert.Equal(t, comment.Body, commentInput.Body)
				assert.Equal(t, comment.ParentPostID, post.ID)
				assert.True(t, comment.ParentCommentID == nil)
				return
			case <-time.After(time.Second * 5):
				t.Error("timeout waiting for comment in subscription")
				return
			}
		}
	}()

	comment, err := resolver.Mutation().CreateComment(ctx, commentInput)
	if err != nil {
		t.Fatal(err.Error())
	}

	wg.Wait()

	assert.Equal(t, *comment.AuthorID, user.ID)
	assert.Equal(t, comment.Body, commentInput.Body)
	assert.Equal(t, comment.ParentPostID, post.ID)
	assert.True(t, comment.ParentCommentID == nil)
}

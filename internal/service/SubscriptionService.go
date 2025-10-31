package service

import (
	"sync"

	"github.com/k0ch3gar/ozon-task/internal/graph/model"
)

type SubscriptionService struct {
	Subs map[string][]chan *model.Comment
	mu   sync.Mutex
}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{
		Subs: make(map[string][]chan *model.Comment),
		mu:   sync.Mutex{},
	}
}

func (ss *SubscriptionService) Unsubscribe(postId string, ch chan *model.Comment) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	var newSubs []chan *model.Comment
	for i := range ss.Subs[postId] {
		if ss.Subs[postId][i] == ch {
			continue
		}

		newSubs = append(newSubs, ss.Subs[postId][i])
	}

	ss.Subs[postId] = newSubs
}

func (ss *SubscriptionService) Subscribe(postId string, ch chan *model.Comment) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.Subs[postId] = append(ss.Subs[postId], ch)
}

func (ss *SubscriptionService) PubComment(postId string, comment *model.Comment) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	chs, ok := ss.Subs[postId]
	if !ok {
		return
	}

	for _, ch := range chs {
		ch <- comment
	}
}

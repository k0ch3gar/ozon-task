package graph

import (
	"github.com/k0ch3gar/ozon-task/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	us *service.UserService
	ps *service.PostService
	cs *service.CommentService
	ss *service.SubscriptionService
}

func NewResolver(
	us *service.UserService,
	ps *service.PostService,
	cs *service.CommentService,
	ss *service.SubscriptionService,
) *Resolver {
	return &Resolver{
		us: us,
		ps: ps,
		cs: cs,
		ss: ss,
	}
}

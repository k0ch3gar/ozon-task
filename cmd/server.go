package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/k0ch3gar/ozon-task/internal/config"
	graph2 "github.com/k0ch3gar/ozon-task/internal/graph"
	config2 "github.com/k0ch3gar/ozon-task/internal/graph/config"
	handler2 "github.com/k0ch3gar/ozon-task/internal/handler"
	"github.com/k0ch3gar/ozon-task/internal/service"
	"github.com/k0ch3gar/ozon-task/internal/storage"
	"go.uber.org/fx"
)

func main() {
	params := config.NewFlagsConfig()

	fx.New(
		fx.Supply(
			params,
		),
		storage.NewStorageModule(params),
		fx.Provide(
			config2.NewResolverConfig,
			service.NewSubscriptionService,
			service.NewUserService,
			service.NewPostService,
			service.NewCommentService,
			graph2.NewResolver,
			handler2.NewGraphQlServer,
		),
		fx.Invoke(func(srv *handler.Server, params config.ApplicationParameters) {
			port := params.Port

			if params.Debug {
				http.Handle("/", playground.Handler("GraphQL playground", "/query"))
				log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
			}

			http.Handle("/query", srv)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		}),
	).Run()
}

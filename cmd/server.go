package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/k0ch3gar/ozon-task/internal/config"
	graph2 "github.com/k0ch3gar/ozon-task/internal/graph"
	config2 "github.com/k0ch3gar/ozon-task/internal/graph/config"
	"github.com/k0ch3gar/ozon-task/internal/service"
	"github.com/k0ch3gar/ozon-task/internal/storage"
	"github.com/vektah/gqlparser/v2/ast"
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
		),
		fx.Provide(
			service.NewUserService,
			graph2.NewResolver,
		),
		fx.Invoke(func(config graph2.Config, params config.ApplicationParameters) {
			port := params.Port

			srv := handler.New(graph2.NewExecutableSchema(config))

			srv.AddTransport(transport.Options{})
			srv.AddTransport(transport.GET{})
			srv.AddTransport(transport.POST{})

			srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

			srv.Use(extension.Introspection{})
			srv.Use(extension.AutomaticPersistedQuery{
				Cache: lru.New[string](100),
			})

			http.Handle("/", playground.Handler("GraphQL playground", "/query"))
			http.Handle("/query", srv)

			log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		}),
	).Run()
}

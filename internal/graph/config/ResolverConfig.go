package config

import graph2 "github.com/k0ch3gar/ozon-task/internal/graph"

func NewResolverConfig(resolver *graph2.Resolver) graph2.Config {
	return graph2.Config{
		Resolvers: resolver,
	}
}

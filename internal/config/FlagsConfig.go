package config

import "flag"

type ApplicationParameters struct {
	PersistentStorageType bool
	Port                  string
	StorageShardsCount    uint64
	PageSize              uint64
	Debug                 bool
}

func NewFlagsConfig() ApplicationParameters {
	var params ApplicationParameters
	flag.Uint64Var(&params.StorageShardsCount, "shards-count", 16, "storage shards count")
	flag.Uint64Var(&params.PageSize, "page-size", 20, "page size")
	flag.StringVar(&params.Port, "port", "8080", "application port")
	flag.BoolVar(&params.PersistentStorageType, "storage-type", true, "storage type where 'true' is persistent storage")
	flag.BoolVar(&params.Debug, "debug", true, "turns on graphql playground")
	flag.Parse()

	return params
}

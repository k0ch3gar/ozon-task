package config

import "flag"

type ApplicationParameters struct {
	StorageType        bool
	Port               string
	StorageShardsCount uint
}

func NewFlagsConfig() ApplicationParameters {
	var params ApplicationParameters
	flag.UintVar(&params.StorageShardsCount, "shards-count", 16, "storage shards count")
	flag.StringVar(&params.Port, "port", "8080", "application port")
	flag.BoolVar(&params.StorageType, "storage-type", true, "storage type where 'true' is persistent storage")
	flag.Parse()

	return params
}

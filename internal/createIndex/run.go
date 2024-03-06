package createIndex

import (
	"go-pentor-bank/internal/config"
)

const (
	prod = "prod"
	sit  = "sit"

	shard   = "shard"
	replica = "replica"
)

var (
	mapPbankDB DBConfig
)

func Run() {

	mapPbankDB = make(map[string]DBConfigDetail)

	// sit
	mapPbankDB[replica] = DBConfigDetail{
		URI:    config.Conf.MongoDriver.DB.URL,
		DBName: config.Conf.MongoDriver.DB.Schema,
	}

	mapPbankDB[shard] = DBConfigDetail{
		URI:    config.Conf.MongoDriver.Shard.URL,
		DBName: config.Conf.MongoDriver.Shard.Schema,
	}

	CreateIndex(mapPbankDB, "pbank")
	//CreateIndex(mapPbankDB, prod, "pbank")
}

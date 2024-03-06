package mongodb

import (
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoDBCon mongoCon

type mongoCon struct {
	DB                *mongo.Client
	SecureClient      SecureClient
	SecureClientShard SecureClient
}

type SecureClient struct {
	MapDataKey   map[string]primitive.Binary
	SecureClient *mongo.Client
	ClientEnc    *mongo.ClientEncryption
	Schema       string
}

func NewMongoDB(log *zerolog.Logger) error {
	secureConn, err := newMongoExplicitEncryptionConn(log, newMongoExpliEncConn{
		DbUrl:     config.Conf.MongoDriver.DB.URL,
		Schema:    config.Conf.MongoDriver.DB.Schema,
		MasterKey: config.Conf.MongoDriver.DB.MasterKey,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("newMongoExplicitEncryptionConn mongoDB got error")
		return err
	}
	MongoDBCon.SecureClient = secureConn

	return nil
}

func NewMongoDBShard(log *zerolog.Logger) error {
	secureConn, err := newMongoExplicitEncryptionConn(log, newMongoExpliEncConn{
		DbUrl:     config.Conf.MongoDriver.Shard.URL,
		Schema:    config.Conf.MongoDriver.Shard.Schema,
		MasterKey: config.Conf.MongoDriver.DB.MasterKey,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("newMongoExplicitEncryptionConn mongoDB Shard got error")
		return err
	}
	MongoDBCon.SecureClientShard = secureConn

	return nil
}

package mongodb

import (
	"context"
	"fmt"
	"go-pentor-bank/internal/config"
	"time"

	"github.com/rs/zerolog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoGoDriverCon(log *zerolog.Logger) error {
	log.Info().Msg("Connecting to mongodb..")

	db, err := DialMongoWithUrl(log, config.Conf.MongoDriver.DB.URL, config.Conf.MongoDriver.DB.Schema)
	if err != nil {
		log.Err(err).Msgf("failed to connect '%v'", config.Conf.MongoDriver.DB.URL)
		return err
	}
	MongoDBCon.DB = db
	MongoDBCon.SecureClient = SecureClient{
		SecureClient: db,
	}

	log.Info().Msg("Connecting to mongodb success!!")
	return nil
}

func DialMongoWithUrl(log *zerolog.Logger, url, dbName string) (*mongo.Client, error) {
	ops := options.Client().ApplyURI(url)
	ops.SetReadPreference(readpref.SecondaryPreferred())
	ops.SetMaxPoolSize(20000)
	ops.SetConnectTimeout(5 * time.Second)
	ops.SetMaxConnIdleTime(10 * time.Second)

	client, err := mongo.Connect(context.Background(), ops)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database %v %v", dbName, err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't ping to database %v %v", dbName, err)
	}
	return client, nil
}

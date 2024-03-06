package createIndex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"slices"
	"time"
)

type KeyIndex []struct {
	CollectionName string  `json:"collectionName"`
	Db             string  `json:"db"`
	DropFirst      bool    `json:"dropFirst"`
	Index          []Index `json:"index"`
}
type Keys struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}
type Options struct {
	Unique      bool `json:"unique"`
	Sparse      bool `json:"sparse"`
	Hidden      bool `json:"hidden"`
	ExpireAfter int  `json:"expireAfter"`
}
type Index struct {
	Keys    []Keys  `json:"keys"`
	Options Options `json:"options"`
}

type DBConfig map[string]DBConfigDetail
type DBConfigDetail struct {
	URI    string
	DBName string
}

func CreateIndex(dbConfig DBConfig, folderName string) {
	log := clog.GetLog()

	mapDbClient := map[string]*mongo.Database{}

	log.Info().Msg("connecting mongoDB~")
	for k, v := range dbConfig {
		msg := fmt.Sprintf("DB : %v", time.Now().Format(time.RFC3339))
		log.Info().Msgf("%s connecting...", msg)
		mongoClient, err := mongodb.DialMongoWithUrl(log, v.URI, v.DBName) //need to modify to v2 ?
		if err != nil {
			log.Error().Err(err).Msg("DB: " + v.DBName)
			msg = fmt.Sprintf("%vcannot dial mongo err: %v", msg, err.Error())
			return
		}

		defer func(mongoClient *mongo.Client, ctx context.Context) {
			err = mongoClient.Disconnect(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("disconnect err: %v", err)
				return
			}
			log.Info().Msg("disconnected..")
		}(mongoClient, context.Background())

		db := mongoClient.Database(v.DBName)
		mapDbClient[k] = db

	}
	log.Info().Msg("Connect to DB successful")

	path := fmt.Sprintf("configs/mongodb/createIndex/index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msgf("ReadFile failed: %v", err)
		return
	}

	var index KeyIndex
	err = json.Unmarshal(data, &index)
	if err != nil {
		log.Error().Err(err).Msgf("Unmarshal failed: %v", err)
		return
	}

	for _, clt := range index {
		var mod []mongo.IndexModel
		for _, keysIndex := range clt.Index {
			keys := bson.D{}
			for _, k := range keysIndex.Keys {
				//keys = keys.Append(k.Key, bsonx.Int64(k.Value))
				keys = append(keys, bson.E{Key: k.Key, Value: k.Value})
			}

			opts := keysIndex.Options
			indexOptions := options.Index().SetUnique(opts.Unique).SetSparse(opts.Sparse).SetHidden(opts.Hidden)

			if len(keysIndex.Keys) > 1 && opts.ExpireAfter > 0 {
				log.Error().Err(errors.New("TTL not support compound index. Please use single key index instead")).Msgf("clt.index %v", clt.CollectionName)
				return
			}

			if opts.ExpireAfter > 0 {
				indexOptions.SetExpireAfterSeconds(int32(opts.ExpireAfter))
			}

			model := mongo.IndexModel{
				Keys:    keys,
				Options: indexOptions,
			}

			mod = append(mod, model)
		}

		collections, err := mapDbClient[clt.Db].ListCollectionNames(context.Background(), bson.M{})
		if slices.Contains(collections, clt.CollectionName) {
			if clt.DropFirst {
				_, err := mapDbClient[clt.Db].Collection(clt.CollectionName).Indexes().DropAll(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("drop index failed: " + clt.CollectionName)
					return
				}
			}

			_, err = mapDbClient[clt.Db].Collection(clt.CollectionName).Indexes().CreateMany(context.Background(), mod)
			if err != nil {
				log.Error().Err(err).Msg("index failed: " + clt.CollectionName)
				return
			}
			log.Info().Msgf("create index for collection: %v total: %v keys is completed", clt.CollectionName, len(clt.Index))
		} else {
			err = mapDbClient[clt.Db].CreateCollection(context.Background(), clt.CollectionName)
			if err != nil {
				log.Error().Err(err).Msg("CreateCollection failed: " + clt.CollectionName)
				return
			}
			log.Info().Msgf("create collection: %v is completed", clt.CollectionName)

			_, err = mapDbClient[clt.Db].Collection(clt.CollectionName).Indexes().CreateMany(context.Background(), mod)
			if err != nil {
				log.Error().Err(err).Msg("index failed: " + clt.CollectionName)
				return
			}
			log.Info().Msgf("create index for collection: %v total: %v keys is completed", clt.CollectionName, len(clt.Index))
		}
	}
}

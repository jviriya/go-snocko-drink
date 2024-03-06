package infra

import (
	"context"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/infra/lcache"
	"go-pentor-bank/internal/infra/mongodb"
	"go-pentor-bank/internal/infra/redis"
	"go-pentor-bank/internal/pkg/aws/s3"
	"go-pentor-bank/internal/pkg/aws/simplequeue"
)

func EstablishInfraConnection(app string) {
	log := clog.GetLog()

	//if err := mongodb.NewMongoGoDriverCon(log); err != nil {
	//	log.Fatal().Err(err).Msg("NewMongoGoDriverCon error")
	//	return
	//}
	//if err := firebaseadmin.InitFirebasePkg(); err != nil {
	//	log.Panic().Err(err).Msg("InitFirebasePkg")
	//	return
	//}

	if err := mongodb.NewMongoDB(log); err != nil {
		log.Fatal().Err(err).Msg("NewMongoDB error")
		return
	}

	if err := mongodb.NewMongoDBShard(log); err != nil {
		log.Fatal().Err(err).Msg("NewMongoDBShard error")
		return
	}

	if err := redis.NewRedisClient(log); err != nil {
		log.Fatal().Err(err).Msg("NewRedisClient error")
		return
	}

	//if err := elasticsearch.NewElasticSearchCon(log); err != nil {
	//	log.Fatal().Err(err).Msg("NewElasticSearchCon error")
	//	return
	//}

	//if err := dynamodb.NewDynamoDBClient(log); err != nil {
	//	log.Fatal().Err(err).Msg("NewDynamoDB error")
	//	return
	//}

	simplequeue.NewAwsSQS()

	s3.NewAwsS3()

	// start local cache storage
	lcache.Init()

}

func ShutdownInfra() {
	//mongodb.MongoDBCon.DB.Disconnect(context.Background())
	if mongodb.MongoDBCon.SecureClient.SecureClient != nil {
		mongodb.MongoDBCon.SecureClient.SecureClient.Disconnect(context.Background())
	}
	if redis.RedisClient != nil {
		redis.RedisClient.Client.Shutdown(context.Background())
	}
}

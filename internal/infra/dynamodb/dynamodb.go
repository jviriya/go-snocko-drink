package dynamodb

import (
	"context"
	configAWS "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/config"
)

var DynamoDBConn dynamodbClient

type dynamodbClient struct {
	Client *dynamodb.Client
}

func NewDynamoDBClient(log *zerolog.Logger) error {
	log.Info().Msg("Connecting to DynamoDB..")

	cfg, err := configAWS.LoadDefaultConfig(
		context.Background(),
		configAWS.WithRegion(config.Conf.AwsConfig.DynamoDB.Region),
		configAWS.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.Conf.AwsConfig.AccessKey,
			config.Conf.AwsConfig.SecretKey,
			""),
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("LoadDefaultConfig got err")
		return err
	}

	DynamoDBConn = dynamodbClient{
		Client: dynamodb.NewFromConfig(cfg),
	}

	log.Info().Msg("Connecting to DynamoDB success!!")

	return nil
}

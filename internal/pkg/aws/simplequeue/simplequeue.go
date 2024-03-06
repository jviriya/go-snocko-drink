package simplequeue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	configAWS "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"strings"
	"sync"
	"time"
)

var (
	sqsClient *sqs.Client
	queueURL  map[string]*string
	initial   sync.Once
)

type (
	Client struct {
		SqsClient *sqs.Client
		queueURL  *string
	}
)

func init() {
	queueURL = make(map[string]*string)
}

func NewAwsSQS() {
	log := clog.GetLog()

	initial.Do(func() {
		cfg, err := configAWS.LoadDefaultConfig(
			context.Background(),
			configAWS.WithRegion(config.Conf.AwsConfig.DefaultRegion),
			configAWS.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				config.Conf.AwsConfig.AccessKey,
				config.Conf.AwsConfig.SecretKey,
				""),
			),
		)
		if err != nil {
			log.Panic().Err(err).Msg("LoadDefaultConfig got err")
		}
		sqsClient = sqs.NewFromConfig(cfg)
	})
}

func getQueueURL(queueName string) (*string, error) {
	queueName = strings.TrimSpace(queueName)
	v, ok := queueURL[queueName]
	if !ok {
		gQInput := &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		}

		result, err := sqsClient.GetQueueUrl(context.TODO(), gQInput)
		if err != nil {
			return nil, err
		}
		queueURL[queueName] = result.QueueUrl
		v = result.QueueUrl
	}
	return v, nil
}

func SendMsg(c context.Context, input *sqs.SendMessageInput, queueName string) (*sqs.SendMessageOutput, error) {
	queueUrl, err := getQueueURL(queueName)
	if err != nil {
		return nil, err
	}
	input.QueueUrl = queueUrl
	return sqsClient.SendMessage(c, input)
}

func GetMessages(c context.Context, input *sqs.ReceiveMessageInput, queueName string) (*sqs.ReceiveMessageOutput, error) {
	queueUrl, err := getQueueURL(queueName)
	if err != nil {
		return nil, err
	}
	input.QueueUrl = queueUrl
	return sqsClient.ReceiveMessage(c, input)
}

func DeleteMessage(c context.Context, input *sqs.DeleteMessageInput, queueName string) (*sqs.DeleteMessageOutput, error) {
	queueUrl, err := getQueueURL(queueName)
	if err != nil {
		return nil, err
	}
	input.QueueUrl = queueUrl
	return sqsClient.DeleteMessage(c, input)
}

func DeleteMessageBatch(c context.Context, input *sqs.DeleteMessageBatchInput, queueName string) (*sqs.DeleteMessageBatchOutput, error) {
	queueUrl, err := getQueueURL(queueName)
	if err != nil {
		return nil, err
	}
	input.QueueUrl = queueUrl
	return sqsClient.DeleteMessageBatch(c, input)
}

func MaxDelaySecond() time.Duration {
	return 15 * time.Minute
}

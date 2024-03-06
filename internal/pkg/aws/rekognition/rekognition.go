package rekognition

import (
	"context"
	configAWS "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"sync"
)

var (
	rekognitionClient *rekognition.Client
	initial           sync.Once
)

func NewAwsRekognition() {
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
		rekognitionClient = rekognition.NewFromConfig(cfg)
	})
}

func CompareFaces(c context.Context, input *rekognition.CompareFacesInput) (*rekognition.CompareFacesOutput, error) {
	return rekognitionClient.CompareFaces(c, input)
}

func CompareFace(c context.Context, source, dest []byte) (float32, error) {
	input := &rekognition.CompareFacesInput{
		//SimilarityThreshold: aws.Float64(90.000000),
		SourceImage: &types.Image{
			Bytes: source,
		},
		TargetImage: &types.Image{
			Bytes: dest,
		},
	}

	var similar float32

	result, err := CompareFaces(c, input)
	if err == nil {
		for _, matchedFace := range result.FaceMatches {
			similar = *matchedFace.Similarity
			break
		}
	}
	return similar, err

}

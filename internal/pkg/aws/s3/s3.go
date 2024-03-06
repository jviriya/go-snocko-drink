package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	configAWS "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type (
	S3FilePath string
)

const (
	KYCLV1Type    S3FilePath = "kyc/level1/"
	KYCLV2Type    S3FilePath = "kyc/level2/"
	ChatType      S3FilePath = "chat/"
	OtherType     S3FilePath = "other/"
	OrderSlipType S3FilePath = "order/slip/"
)

var (
	s3Client *s3.Client
	initial  sync.Once
)

func NewAwsS3() {
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
			log.Error().Err(err).Msg("LoadDefaultConfig got err")
		}
		s3Client = s3.NewFromConfig(cfg)
	})
}

func ListBuckets(log *zerolog.Logger) ([]types.Bucket, error) {
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	var buckets []types.Bucket
	if err != nil {
		log.Error().Err(err).Msg("Couldn't list buckets for your account. Here's why: ")
	} else {
		buckets = result.Buckets
	}
	return buckets, err
}

func CreateBucket(log *zerolog.Logger, name string) error {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(config.Conf.AwsConfig.DefaultRegion),
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't create bucket. Here's why: ")
	}
	return err
}

func DeleteBucket(log *zerolog.Logger, bucketName string) error {
	_, err := s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't delete bucket. Here's why: ")
	}
	return err
}

func DeleteObjects(log *zerolog.Logger, bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := s3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't delete objects from bucket. Here's why: ")
	}
	return err
}

func DownloadFile(log *zerolog.Logger, path string) ([]byte, error) {
	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(config.Conf.AwsConfig.S3.CDN),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't get object. Here's why: ")
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't read object body from. Here's why: %v")
	}

	return body, err
}

func ListObjects(log *zerolog.Logger, bucketName string) ([]types.Object, error) {
	result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []types.Object
	if err != nil {
		log.Error().Err(err).Msg("Couldn't list objects in bucket. Here's why: ")
	} else {
		contents = result.Contents
	}
	return contents, err
}

func UploadFile(log *zerolog.Logger, fileType S3FilePath, fromFile *multipart.FileHeader) (string, error) {
	f, err := fromFile.Open()

	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := fmt.Sprintf("%02d", int(t.Month()))

	objectKey := string(fileType) + y + m + "/" + uuid.New().String() + filepath.Ext(fromFile.Filename)

	if err != nil {
		log.Error().Err(err).Msg("Couldn't open file to upload. Here's why: ")
	} else {
		if fromFile.Size > config.Conf.Image.MaxFileSize {
			log.Error().Err(errors.New(common.ErrFileSizeLimitation)).Msgf("size exceed the limitation got")
			return "", errors.New(common.ErrFileSizeLimitation)
		}

		if _, found := config.Conf.Image.AllowedMimeTypesMap[fromFile.Header.Get("Content-Type")]; !found {
			log.Error().Err(errors.New(common.ErrInvalidFileMIME)).Msg("invalid f mime")
			return "", errors.New(common.ErrInvalidFileMIME)
		}

		if _, found := config.Conf.Image.AllowedExtensionsMap[filepath.Ext(fromFile.Filename)]; !found {
			log.Error().Err(errors.New(common.ErrInvalidFileExtension)).Msg("invalid file extension")
			return "", errors.New(common.ErrInvalidFileExtension)
		}

		defer f.Close()
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(config.Conf.AwsConfig.S3.CDN),
			Key:         aws.String(objectKey),
			Body:        f,
			ContentType: aws.String(fromFile.Header.Get("Content-Type")),
		})
		if err != nil {
			log.Error().Err(err).Msg("Couldn't upload file. Here's why: ")
		}
	}
	return "https://" + config.Conf.AwsConfig.S3.CDN + "/" + objectKey, err
}

func UploadFileFormReader(log *zerolog.Logger, fileType S3FilePath, file io.Reader) (string, error) {
	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := fmt.Sprintf("%02d", int(t.Month()))

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	contentType := http.DetectContentType(buf.Bytes())
	fileExtension := ".jpg"
	if contentType == "image/png" {
		fileExtension = ".png"
	} else if contentType == "image/webp" {
		fileExtension = ".webp"
	}
	objectKey := string(fileType) + y + m + "/" + uuid.New().String() + fileExtension

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.Conf.AwsConfig.S3.CDN),
		Key:         aws.String(objectKey),
		Body:        buf,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't upload file. Here's why: ")
	}
	return "https://" + config.Conf.AwsConfig.S3.CDN + "/" + objectKey, err
}

func CreateFolder(log *zerolog.Logger, bucketName, objectKey string) error {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Error().Err(err).Msg("Couldn't upload file. Here's why: ")
	}
	return err
}

package images

import (
	"bytes"
	"encoding/base64"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/pkg/aws/s3"
	"go-pentor-bank/internal/utils"
	"strings"
)

func UploadResizeToS3(base64String string, log *zerolog.Logger) (string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))

	b, err := utils.ResizeImage(reader)
	if err != nil {
		log.Error().Err(err).Msg("utils.ResizeImage got err")
		return "", err
	}

	s3Upload := bytes.NewReader(b)

	url, err := s3.UploadFileFormReader(log, s3.ChatType, s3Upload)
	if err != nil {
		log.Error().Err(err).Msg("UploadImage got err")
		return "", err
	}

	return url, nil
}

func UploadNormalSizeToS3(base64String string, log *zerolog.Logger) (string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))

	b, err := utils.ResizeNormalImage(reader)
	if err != nil {
		log.Error().Err(err).Msg("utils.ResizeImage got err")
		return "", err
	}

	s3Upload := bytes.NewReader(b)

	url, err := s3.UploadFileFormReader(log, s3.ChatType, s3Upload)
	if err != nil {
		log.Error().Err(err).Msg("UploadImage got err")
		return "", err
	}

	return url, nil
}

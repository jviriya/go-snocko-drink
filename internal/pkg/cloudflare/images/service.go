package images

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/curl"
	"go-pentor-bank/internal/utils"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	TypeKYC        ImageType = "KYC"
	TypeChat       ImageType = "CHAT"
	TypeToken      ImageType = "TOKEN"
	TypeCurrency   ImageType = "CURRENCY"
	TypeCountry    ImageType = "COUNTRY"
	TypeProfilePic ImageType = "PROFILE_PICTURE"
)

type (
	ImageType string
	Metadata  struct {
		Type     ImageType `json:"type"`
		Filename string    `json:"filename"`
	}
)

func UploadImage(log *zerolog.Logger, file io.Reader, imgType ImageType, filename string) ([]string, error) {
	meta := Metadata{
		Type:     imgType,
		Filename: filename,
	}
	metaB, _ := sonic.Marshal(meta)
	//mp := multipart.FormData{
	//	Files: []multipart.FormFile{
	//		{
	//			Name:   "file",
	//			Reader: file,
	//		},
	//	},
	//	Data: map[string]multipart.Values{
	//		"metadata": {string(metaB)},
	//	},
	//}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	metaWriter, _ := writer.CreateFormField("metadata")
	metaWriter.Write(metaB)
	fileWriter, _ := writer.CreateFormFile("file", filename)
	io.Copy(fileWriter, file)
	writer.Close()

	strBuilder := strings.NewReplacer(":accountID", config.Conf.CloudFlare.Images.AccountID)
	path := strBuilder.Replace(config.Conf.CloudFlare.Images.Path.Upload)

	url := fmt.Sprintf("%s/%s", config.Conf.CloudFlare.Domain, path)
	body, _, _, err := curl.NewMultipartData(url, buff.Bytes(), nil, writer).Bearer(config.Conf.CloudFlare.Images.ApiToken).POST()
	if err != nil {
		log.Error().Err(err).Msg("curl.NewRest ")
		return nil, err
	}

	var data UploadImageResp
	err = sonic.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if !data.Success {
		log.Error().Msgf("found %v error(s)", len(data.Errors))
		for _, v := range data.Errors {
			log.Error().Msgf("code %v, msg: %s", v.Code, v.Message)
		}
		return nil, errors.New(common.ErrInternalError)
	}

	var variants []string
	for _, v := range data.Result.Variants {
		if strings.HasSuffix(v, "public") {
			variants = append(variants, v)
		}
	}

	return variants, nil
}

func GetImage(log *zerolog.Logger) ([]ListImages, error) {

	header := utils.GenerateHeaderJson()

	strBuilder := strings.NewReplacer(":accountID", config.Conf.CloudFlare.Images.AccountID)
	path := strBuilder.Replace(config.Conf.CloudFlare.Images.Path.ListImage)

	url := fmt.Sprintf("%s/%s", config.Conf.CloudFlare.Domain, path)
	body, statusCode, _, err := curl.NewRest(url, nil, header).Bearer(config.Conf.CloudFlare.Images.ApiToken).GET()
	if err != nil {
		log.Error().Err(err).Msg("curl.NewRest ")
		return nil, err
	}

	if statusCode != http.StatusOK {
		log.Error().Msg("status is not ok")
		log.Error().Msg(fmt.Sprint(body))
		return nil, err
	}

	var data []ListImages
	err = sonic.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

package service

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/infra/lcache"
	"go-pentor-bank/internal/pkg/curl"
	"net/http"
)

type ImgCacheContent struct {
	Content []byte `json:"content"`
	Header  string `json:"header"`
}

func GetImageForStore(main, sub, url string) ([]byte, string, error) {
	var data ImgCacheContent

	body, statusCode, headerResp, err := curl.NewRest(url, nil, nil).GET()
	if err != nil {
		log.Error().Err(err).Msg("curl.NewRest ")
		return nil, "", nil
	}

	if statusCode != http.StatusOK {
		log.Error().Msg("status is not ok")
		log.Error().Msg(string(body))
		return nil, "", nil
	}

	//b, err = io.ReadAll(body)
	//if err != nil {
	//	log.Error().Msg("io.ReadAll got err")
	//	return nil, "", err
	//}
	hd := headerResp.Get("Content-Type")

	data = ImgCacheContent{
		Content: body,
		Header:  hd,
	}

	key := fmt.Sprintf("%s_%s_%s", common.CdnImgCache, main, sub)
	err = lcache.Set(key, data, common.CdnImgCacheExpire)
	if err != nil {
		log.Error().Msg("lcache.Set got err")
		return nil, "", err
	}

	return body, hd, nil

}

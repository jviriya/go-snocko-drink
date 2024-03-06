package geetest

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	RequestStatusSuccess = "success"
	RequestStatusError   = "error"

	VerificationResultSuccess = "success"
	VerificationResultFail    = "fail"
)

func CallbackValidate(log *zerolog.Logger, URL, lotNumber, captchaOutput, passToken, genTime, signToken string) error {
	formData := make(url.Values)
	formData["lot_number"] = []string{lotNumber}
	formData["captcha_output"] = []string{captchaOutput}
	formData["pass_token"] = []string{passToken}
	formData["gen_time"] = []string{genTime}
	formData["sign_token"] = []string{signToken}

	cli := http.Client{Timeout: time.Second * 5}
	resp, err := cli.PostForm(URL, formData)
	if err != nil || resp.StatusCode != 200 {
		log.Error().Err(err).Msg("cli.PostForm got err")
		return err
	}

	respJson, _ := ioutil.ReadAll(resp.Body)
	var respMap map[string]interface{}
	if err = json.Unmarshal(respJson, &respMap); err != nil {
		log.Error().Err(err).Msg("json.Unmarshal got err")
		return err
	}

	switch respMap["status"] {
	case RequestStatusSuccess:
		if respMap["result"] != VerificationResultSuccess {
			return errors.New(common.ErrInternalError)
		}
	default:
		return errors.New(common.ErrInternalError)
	}

	return nil
}

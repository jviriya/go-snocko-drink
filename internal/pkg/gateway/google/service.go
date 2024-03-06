package google

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/pkg/curl"
	"go-pentor-bank/internal/utils"
	"time"
)

type ValidationResp struct {
	Success        bool      `json:"success"`
	ChallengeTS    time.Time `json:"challenge_ts"`
	Hostname       string    `json:"hostname"`
	ApkPackageName string    `json:"apkPackageName"`
	ErrorCodes     []string  `json:"error-codes"`
}

func CallbackVerifyReCaptcha(log *zerolog.Logger, URL, secretKey, responseKey string) error {

	headers := utils.GenerateHeaderJson()
	pathURL := fmt.Sprintf("%s?secret=%s&response=%s", URL, secretKey, responseKey)

	body, _, _, err := curl.NewRest(pathURL, nil, headers).POST()
	if err != nil {
		log.Error().Err(err).Msgf("curl.NewRest got err")
		return errors.New(common.ErrInternalError)
	}

	var resp ValidationResp
	err = sonic.Unmarshal(body, &resp)
	if err != nil {
		log.Error().Err(err).Msg("read response from request got err")
		return errors.New(common.ErrInternalError)
	}

	if !resp.Success {
		errNew := errors.New(fmt.Sprintf("response got err code %v", resp.ErrorCodes))
		log.Error().Err(errNew).Msg("turnstile response got err")
		return errNew
	}

	return nil
}

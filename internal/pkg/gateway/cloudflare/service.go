package cloudflare

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

type ValidationForm struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	RemoteIP string `json:"remoteip"`
}

type ValidationResp struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
	Action      string    `json:"action"`
	Cdata       string    `json:"cdata"`
}

func CallbackTurnstile(log *zerolog.Logger, URL, secretKey, responseKey, remoteIP string) error {

	data := ValidationForm{
		Secret:   secretKey,
		Response: responseKey,
		RemoteIP: remoteIP,
	}

	headers := utils.GenerateHeaderJson()

	body, _, _, err := curl.NewRest(URL, data, headers).POST()
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

package vonage

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/curl"
	"go-pentor-bank/internal/utils"
	"strconv"
)

type vonage struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	From      string `json:"from"`
	To        string `json:"to"`
	Text      string `json:"text"`
}

func SendSMS(log *zerolog.Logger, mobileNumber, msg string) error {
	endpoint := config.Conf.Vonage.API + config.Conf.Vonage.Path.SendSMS
	headers := utils.GenerateHeaderJson()

	data := vonage{
		ApiKey:    config.Conf.Vonage.APIKey,
		ApiSecret: config.Conf.Vonage.APISecret,
		From:      config.Conf.Vonage.From,
		To:        mobileNumber,
		Text:      msg,
	}

	body, _, _, err := curl.NewRest(endpoint+"?", data, headers).POST()
	if err != nil {
		log.Error().Err(err).Msg("curl sms got err")
		return err
	}

	var resp SendSMSResp
	err = sonic.Unmarshal(body, &resp)
	if err != nil {
		log.Error().Err(err).Msg("read response from request got err")
		return err
	}

	count, err := strconv.ParseInt(resp.MessageCount, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("ParseInt msg count got err")
		return err
	}

	if count <= 0 {
		return errors.New(ErrNoSentSMS)
	}

	if resp.Messages[0].Status != "0" {
		return errors.New(ErrSendingFailed)
	}

	return nil
}

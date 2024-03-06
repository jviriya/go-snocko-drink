package deecommerce

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/curl"
)

func SendSMS(to, msg string) error {

	req := SendSMSReq{
		AccountID: config.Conf.SMS.Deecommerce.SmsAccountID,
		SecretKey: config.Conf.SMS.Deecommerce.SmsSecretKey,
		Type:      config.Conf.SMS.Deecommerce.SmsType,
		To:        to,
		Sender:    config.Conf.SMS.Deecommerce.SmsSender,
		Msg:       msg,
	}

	headers := curl.GenerateHeaderJson()

	urlRequest := fmt.Sprintf("%s/service/SMSWebService", config.Conf.SMS.Deecommerce.Domain)

	body, _, _, err := curl.NewRest(urlRequest, req, headers).POST()
	if err != nil {
		return err
	}

	rtn := SendSMSResp{}
	err = sonic.Unmarshal(body, &rtn)
	if err != nil {
		return err
	}

	if rtn.Error != "0" && rtn.Msg == "Invalid telephone format." {
		return nil
	} else if rtn.Error != "0" && rtn.Msg == "Your message is blocked." {
		return nil
	} else if rtn.Error != "0" {
		return errors.New(rtn.Msg)
	}

	return nil
}

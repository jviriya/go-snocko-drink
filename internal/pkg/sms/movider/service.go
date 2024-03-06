package movider

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/curl"
	"net/http"
	"net/url"
)

func GetBalance(c *fiber.Ctx) error {

	log := clog.GetContextLog(c.Context())
	data := url.Values{}
	data.Set("api_key", config.Conf.SMS.Movider.ApiKey)
	data.Set("api_secret", config.Conf.SMS.Movider.ApiSecret)

	headers := curl.GenerateHeaderUrlencoded()

	headers["accept"] = "application/json"
	headers["content-type"] = "application/x-www-form-urlencoded"

	resp, statusCode, _, err := curl.NewRest(fmt.Sprintf("%s%s", config.Conf.SMS.Movider.Domain, getBalance), data.Encode(), headers).POST()
	if err != nil {
		log.Error().Err(err).Msg("curl.NewRest ")
		return err
	}

	if statusCode != http.StatusOK {
		log.Error().Msgf("body %v", string(resp))
		return errors.New(fmt.Sprint(statusCode))
	}

	return nil
}

func SendSMS(c context.Context, tel, text string) error {

	log := clog.GetContextLog(c)

	data := url.Values{}
	data.Set("api_key", config.Conf.SMS.Movider.ApiKey)
	data.Set("api_secret", config.Conf.SMS.Movider.ApiSecret)
	data.Set("text", text)
	data.Set("to", tel)

	headers := curl.GenerateHeaderUrlencoded()

	headers["accept"] = "application/json"
	headers["content-type"] = "application/x-www-form-urlencoded"

	body, statusCode, _, err := curl.NewRest(fmt.Sprintf("%s%s", config.Conf.SMS.Movider.Domain, sendSMS), data.Encode(), headers).POST()
	if err != nil {
		log.Error().Err(err).Msg("curl.NewRest ")
		return err
	}

	if statusCode != http.StatusOK {
		log.Error().Msgf("got error, body : %v", string(body))
		return errors.New(fmt.Sprint(statusCode))
	}

	var resp Response
	err = sonic.Unmarshal(body, &resp)
	if err != nil {
		log.Error().Err(err).Msg("sonic.Unmarshal got err")
		return err
	}

	if resp.Error.Code != 0 {
		log.Error().Msgf("error code is not equal zero, body : %v", string(body))
		return errors.New("error code is not equal zero")
	}

	return nil
}

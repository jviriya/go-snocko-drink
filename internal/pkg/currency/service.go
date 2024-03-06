package currency

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/curl"
	"go-pentor-bank/internal/utils"
	"net/http"
)

func GetCurrencyRates(base, symbols string) (RatesResp, error) {

	log := clog.GetLog()
	resp := RatesResp{}

	endPoint := fmt.Sprintf("%s?apikey=%s&currencies=%s&base_currency=%s",
		config.Conf.CurrencyApi.CurrencyApiUrl,
		config.Conf.CurrencyApi.CurrencyAccessKey,
		symbols,
		base)

	headers := utils.GenerateHeaderJson()

	body, statusCode, _, err := curl.NewRest(endPoint, nil, headers).GET()
	if err != nil {
		log.Error().Err(err).Msg("curl.NewRest err: ")
		return resp, err
	}

	if statusCode != http.StatusOK {
		errNew := errors.New(fmt.Sprintf("status != 200, got %v", statusCode))
		return resp, errNew
	}

	err = sonic.Unmarshal(body, &resp)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal got err, ")
		return resp, err
	}

	return resp, nil
}

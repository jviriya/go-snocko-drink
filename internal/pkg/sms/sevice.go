package sms

import (
	"context"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/pkg/sms/deecommerce"
	"go-pentor-bank/internal/pkg/sms/movider"
	"go-pentor-bank/internal/pkg/sms/vonage"
	"strings"
)

func SendSMS(c context.Context, log *zerolog.Logger, req SendSMSReq) error {
	var err error
	if strings.HasPrefix(req.Tel, "66") {
		err = deecommerce.SendSMS(req.Tel, req.Text)
		if err != nil {
			log.Error().Err(err).Msg("deecommerce.SendSMS got err")
			return err
		}
	} else if strings.HasPrefix(req.Tel, "855") {
		err = vonage.SendSMS(log, req.Tel, req.Text)
		if err != nil {
			log.Error().Err(err).Msg("dvonage.SendSMS got err")
			return err
		}
	} else {
		err = movider.SendSMS(c, "+"+req.Tel, req.Text)
		if err != nil {
			log.Error().Err(err).Msg("movider.SendSMS got err")
			return err
		}
	}

	return err
}

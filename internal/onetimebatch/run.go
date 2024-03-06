package onetimebatch

import (
	"go-pentor-bank/internal/clog"
)

func Run(job string) {
	log := clog.GetLog()
	switch job {
	case "ui_text":
		//uiText.Run(commondb.IosPlatformUiText, []string{"en", "th"})
		//uiText.Run(commondb.AndroidPlatformUiText, []string{"en", "th"})
	default:
		log.Error().Msgf("job %v is not found!!!", job)
		return
	}
	log.Info().Msgf("running job %v is success!", job)
}

package middleware

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/pkg/jwtManager"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog"
)

var unlogPath = []string{"healthz", "status", "loadtest"}
var hideRequestBodyPath = []string{"login"}

const (
	hideRespBodyKey = "HIDE_RESP_BODY_HIDE"
)

func LogMetadataReqResp(service string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := clog.GetContextLog(c.Context())

		sp := strings.Split(c.Path(), "/")
		if slices.Contains(unlogPath, sp[len(sp)-1]) {
			return c.Next()
		}

		var start, stop time.Time
		start = time.Now()

		var hideReqBodyKey bool
		if slices.Contains(hideRequestBodyPath, sp[len(sp)-1]) {
			hideReqBodyKey = true
		}

		var metadataReq string
		if hideReqBodyKey {
			metadataReq = fmt.Sprintf("FIBER - %s %s", c.Method(), c.Path())
		} else {
			metadataReq = fmt.Sprintf("FIBER - %s %s - reqBody:%s reqString:%s", c.Method(), c.Path(), strings.Replace(string(c.Body()), "\n", "", -1), c.Request().URI().QueryArgs().String())
			if strings.Contains(c.Get("Content-Type"), "multipart/form-data") {
				metadataReq = fmt.Sprintf("FIBER - %s %s -  reqString:%s", c.Method(), c.Path(), c.Request().URI().QueryArgs().String())
			}
		}
		log.Info().Msg(metadataReq)

		ctxErr := c.Next()

		log = clog.GetContextLog(c.Context())

		stop = time.Now()

		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("duration", stop.Sub(start).String())
		})

		code := c.Response().StatusCode()
		//message := c.Response().StatusCode()
		if ctxErr != nil {
			switch val := ctxErr.(type) {
			case *fiber.Error:
				code = ctxErr.(*fiber.Error).Code
			case *json.MarshalerError:
				log.Error().Msgf("%v", val)
			}
		}

		var userLogin jwtManager.SessionData

		if userLoginRaw := c.Locals("userLogin"); userLoginRaw != nil {
			userLogin = userLoginRaw.(jwtManager.SessionData)
		}

		hideRespBody := c.Locals(hideRespBodyKey)
		var metadataResp string
		if hideRespBody != nil && hideRespBody.(bool) {
			metadataResp = fmt.Sprintf("FIBER - %s %s - %d %s uid: %s ", c.Method(), c.Path(), code, userLogin.UID, stop.Sub(start).String())
		} else {
			rsData := strings.Replace(string(c.Response().Body()), "\n", "", -1)
			if utf8.RuneCount([]byte(rsData)) > 300 {
				rsData = string([]rune(rsData)[:300]) + "...too many"
			}
			metadataResp = fmt.Sprintf("FIBER - %s %s - %d %s - uid: %s respBody:%s", c.Method(), c.Path(), code, stop.Sub(start).String(), userLogin.UID, rsData)
		}

		if code == fiber.StatusOK {
			log.Info().Msg(metadataResp)
		} else { // Status not OK
			log.Error().Msg(metadataResp)
		}

		return ctxErr
	}
}

func HideBodyResp(c *fiber.Ctx) error {
	c.Locals(hideRespBodyKey, true)
	return c.Next()
}

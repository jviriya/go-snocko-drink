package middleware

//
//import (
//	"crypto/hmac"
//	"crypto/sha256"
//	"encoding/hex"
//	"errors"
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"github.com/rs/zerolog"
//	"go-pentor-bank/internal/app/auth/request"
//	"go-pentor-bank/internal/appresponse"
//	"go-pentor-bank/internal/clog"
//	"go-pentor-bank/internal/common"
//	"go-pentor-bank/internal/config"
//	"go-pentor-bank/internal/pkg/gateway/geetest"
//	reporedis "go-pentor-bank/internal/repository/redis"
//	"go-pentor-bank/internal/utils"
//	"net/http"
//	"strings"
//)
//
//const (
//	RequestStatusSuccess = "success"
//	RequestStatusError   = "error"
//
//	VerificationResultSuccess = "success"
//	VerificationResultFail    = "fail"
//)
//
//func ValidationCaptchaGeeTest() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		log := clog.GetContextLog(c.Context())
//
//		var captchaID, captchaKey string
//
//		platform := c.Get(common.PlatformContext)
//
//		switch platform {
//		case common.IosPlatform:
//			captchaKey = config.Conf.GeeTestCaptchaV4.IOSCaptchaKey
//
//		case common.AndroidPlatform:
//			captchaKey = config.Conf.GeeTestCaptchaV4.AndroidCaptchaKey
//
//		case common.WebPlatform:
//			captchaKey = config.Conf.GeeTestCaptchaV4.WebCaptchaKey
//
//		default:
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//
//		apiServer := config.Conf.GeeTestCaptchaV4.Api
//		captchaID = c.Get("Captcha-ID")
//		//if strings.Trim(ca)
//		URL := fmt.Sprintf("%s/validate?captcha_id=%s", apiServer, captchaID)
//
//		lotNumber := c.Get("Lot-Number")
//		captchaOutput := c.Get("Captcha-Output")
//		passToken := c.Get("Pass-Token")
//		genTime := c.Get("Gen-Time")
//		signToken := hmac_encode(captchaKey, lotNumber)
//
//		err := geetest.CallbackValidate(log, URL, lotNumber, captchaOutput, passToken, genTime, signToken)
//		if err != nil {
//			log.Error().Err(err).Msg("CallbackValidate got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//
//		return c.Next()
//	}
//}
//
//func ValidationCaptchaLogin() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		log := clog.GetContextLog(c.Context())
//
//		req := request.Login{}
//		err := c.BodyParser(&req)
//		if err != nil {
//			log.Error().Err(err).Msg("CallbackValidate got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//
//		platform := c.Get(common.PlatformContext)
//		switch platform {
//		case common.IosPlatform, common.AndroidPlatform:
//			platform = common.MobilePlatform
//		case common.WebPlatform:
//			platform = common.WebPlatform
//		default:
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//
//		var username string
//		switch req.Type {
//		case common.EmailType:
//			username = req.TypeVal
//		case common.PhoneNoType:
//			phoneNo := utils.GetOnlyPhoneNumber(req.TypeVal)
//			username = req.DialCode + phoneNo
//		default:
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//
//		var key reporedis.CaptchaModule
//		if platform == common.WebPlatform {
//			key = common.CaptchaLoginWebKey
//		} else {
//			key = common.CaptchaLoginMobileKey
//		}
//		status, errCode, err := reporedis.GetCaptcha(c.Context(), log, key, utils.GetIP(c), username)
//		if err != nil {
//			log.Error().Err(err).Msg("GetCaptcha got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: errCode})
//		}
//
//		errCode, err = reporedis.DelCaptcha(c.Context(), log, key, utils.GetIP(c), username)
//		if err != nil {
//			log.Error().Err(err).Msg("DelCaptcha got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: errCode})
//		}
//
//		if status == common.StatusRequired {
//			err = CallbackCaptchaValidate(c, log)
//			if err != nil {
//				log.Error().Err(err).Msg("CallbackCaptchaValidate got err")
//				return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//			}
//		}
//
//		return c.Next()
//	}
//}
//
//func ValidationCaptchaRegister() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		log := clog.GetContextLog(c.Context())
//
//		status, errCode, err := reporedis.GetCaptcha(c.Context(), log, common.CaptchaRegisterKey, utils.GetIP(c), "")
//		if err != nil {
//			log.Error().Err(err).Msg("GetCaptcha got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: errCode})
//		}
//
//		errCode, err = reporedis.DelCaptcha(c.Context(), log, common.CaptchaRegisterKey, utils.GetIP(c), "")
//		if err != nil {
//			log.Error().Err(err).Msg("DelCaptcha got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: errCode})
//		}
//
//		if status == common.StatusRequired {
//			err = CallbackCaptchaValidate(c, log)
//			if err != nil {
//				log.Error().Err(err).Msg("CallbackCaptchaValidate got err")
//				return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//			}
//		}
//
//		return c.Next()
//	}
//}
//
//func ValidationCaptchaResetPassword() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		log := clog.GetContextLog(c.Context())
//
//		status, errCode, err := reporedis.GetCaptcha(c.Context(), log, common.CaptchaResetPasswordKey, utils.GetIP(c), "")
//		if err != nil {
//			log.Error().Err(err).Msg("GetCaptcha got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: errCode})
//		}
//
//		errCode, err = reporedis.DelCaptcha(c.Context(), log, common.CaptchaResetPasswordKey, utils.GetIP(c), "")
//		if err != nil {
//			log.Error().Err(err).Msg("DelCaptcha got err")
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: errCode})
//		}
//
//		if status == common.StatusRequired {
//			err = CallbackCaptchaValidate(c, log)
//			if err != nil {
//				log.Error().Err(err).Msg("CallbackCaptchaValidate got err")
//				return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//			}
//		}
//
//		return c.Next()
//	}
//}
//
//func hmac_encode(key string, data string) string {
//	mac := hmac.New(sha256.New, []byte(key))
//	mac.Write([]byte(data))
//	return hex.EncodeToString(mac.Sum(nil))
//}
//
//func CallbackCaptchaValidate(c *fiber.Ctx, log *zerolog.Logger) error {
//	var captchaID, captchaKey string
//	platform := c.Get(common.PlatformContext)
//	switch platform {
//	case common.WebPlatform:
//		captchaKey = config.Conf.GeeTestCaptchaV4.WebCaptchaKey
//	default:
//		return errors.New("invalid platform")
//	}
//	apiServer := config.Conf.GeeTestCaptchaV4.Api
//	captchaID = c.Get("Captcha-ID")
//	if strings.TrimSpace(captchaID) == "" {
//		log.Error().Err(errors.New("")).Msg("captchaID invalid")
//		return errors.New("CallbackValidate got err")
//	}
//	URL := fmt.Sprintf("%s/validate?captcha_id=%s", apiServer, captchaID)
//	lotNumber := c.Get("Lot-Number")
//	if strings.TrimSpace(lotNumber) == "" {
//		log.Error().Err(errors.New("")).Msg("lotNumber invalid")
//		return errors.New("CallbackValidate got err")
//	}
//	captchaOutput := c.Get("Captcha-Output")
//	if strings.TrimSpace(captchaOutput) == "" {
//		log.Error().Err(errors.New("")).Msg("captchaOutput invalid")
//		return errors.New("CallbackValidate got err")
//	}
//	passToken := c.Get("Pass-Token")
//	if strings.TrimSpace(passToken) == "" {
//		log.Error().Err(errors.New("")).Msg("passToken invalid")
//		return errors.New("CallbackValidate got err")
//	}
//	genTime := c.Get("Gen-Time")
//	if strings.TrimSpace(genTime) == "" {
//		log.Error().Err(errors.New("")).Msg("genTime invalid")
//		return errors.New("CallbackValidate got err")
//	}
//	signToken := hmac_encode(captchaKey, lotNumber)
//
//	err := geetest.CallbackValidate(log, URL, lotNumber, captchaOutput, passToken, genTime, signToken)
//	if err != nil {
//		log.Error().Err(err).Msg("CallbackValidate got err")
//		return errors.New("CallbackValidate got err")
//	}
//
//	return nil
//}

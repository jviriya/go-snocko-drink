package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/jwtManager"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwt *jwtManager.JWTManager
}

func NewAuthMiddleware(jwt *jwtManager.JWTManager) AuthMiddleware {
	return AuthMiddleware{
		jwt: jwt,
	}
}

func (m AuthMiddleware) RetrieveSession(c *fiber.Ctx) error {
	log := clog.GetContextLog(c.Context())

	if c.Get("Authorization") == "" {
		return c.Next()
	}

	platform := c.Get(common.PlatformContext)
	switch strings.TrimSpace(platform) {
	case common.AndroidPlatform, common.IosPlatform:
		platform = common.MobilePlatform
	case common.WebPlatform:
	default:
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.InvalidPlatform})
	}
	c.Locals(common.PlatformContext, platform)

	tokenAuth, err := m.jwt.ExtractTokenMetadata(c.Context(), c.Get("Authorization"))
	if err != nil {
		log.Error().Err(err).Msg("ExtractTokenMetadata :")
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}
	userLogin, err := m.jwt.FetchAuth(c.Context(), tokenAuth.AccessUUID)
	if err != nil {
		log.Error().Err(err).Msg("FetchAuth :")
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}
	c.Locals("userLogin", userLogin)
	c.Locals("uuid", tokenAuth.AccessUUID)
	valid, err := m.jwt.CheckLastAccessUUID(c.Context(), tokenAuth.AccessUUID, tokenAuth.UserID)
	if err != nil {
		log.Error().Err(err).Msg("CheckAccessUUID :")
	}

	if !valid {
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}

	if userLogin.UserStatus != common.StatusActive {
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.InactiveAccount})
	}
	return c.Next()
}

func (m AuthMiddleware) VerifySession(c *fiber.Ctx) error {
	log := clog.GetContextLog(c.Context())

	platform := c.Get(common.PlatformContext)
	switch strings.TrimSpace(platform) {
	case common.WebPlatform:
	case common.IosPlatform, common.AndroidPlatform:
		platform = common.MobilePlatform
	default:
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.InvalidPlatform})
	}
	c.Locals(common.PlatformContext, platform)

	tokenAuth, err := m.jwt.ExtractTokenMetadata(c.Context(), c.Get("Authorization"))
	if err != nil {
		log.Error().Err(err).Msg("ExtractTokenMetadata :")
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}
	userLogin, err := m.jwt.FetchAuth(c.Context(), tokenAuth.AccessUUID)
	if err != nil {
		log.Error().Err(err).Msg("FetchAuth :")
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}

	logger := clog.GetContextLog(c.Context())
	tempLog := logger.With().Fields(map[string]interface{}{
		"uid": userLogin.UID,
	}).Logger()
	*logger = tempLog
	c.Locals(clog.CLoggerKey, logger)

	c.Locals("userLogin", userLogin)
	c.Locals("uuid", tokenAuth.AccessUUID)
	valid, err := m.jwt.CheckLastAccessUUID(c.Context(), tokenAuth.AccessUUID, tokenAuth.UserID)
	if err != nil {
		log.Error().Err(err).Msg("CheckAccessUUID :")
	}

	if !valid {
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}

	if userLogin.UserStatus != common.StatusActive {
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.InactiveAccount})
	}

	return c.Next()
}

func (m AuthMiddleware) VerifyBOSession(c *fiber.Ctx) error {
	log := clog.GetContextLog(c.Context())

	c.Locals(common.PlatformContext, common.BOPlatform)

	tokenAuth, err := m.jwt.ExtractTokenMetadata(c.Context(), c.Get("Authorization"))
	if err != nil {
		log.Error().Err(err).Msg("ExtractTokenMetadata :")
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}
	userLogin, err := m.jwt.FetchAuth(c.Context(), tokenAuth.AccessUUID)
	if err != nil {
		log.Error().Err(err).Msg("FetchAuth :")
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}
	c.Locals("userLogin", userLogin)
	c.Locals("uuid", tokenAuth.AccessUUID)
	valid, err := m.jwt.CheckLastAccessUUID(c.Context(), tokenAuth.AccessUUID, tokenAuth.UserID)
	if err != nil {
		log.Error().Err(err).Msg("CheckAccessUUID :")
	}

	if !valid {
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
	}
	return c.Next()
}

func (m AuthMiddleware) VerifyRole(values ...common.UserType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		mapRole := make(map[common.UserType]struct{})
		for _, v := range values {
			mapRole[v] = struct{}{}
		}

		if len(values) != 0 {
			session := jwtManager.GetSessionData(c)
			if _, ok := mapRole[session.Role]; !ok {
				return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.Permission})
			}

		}

		return c.Next()
	}
}

//func (m AuthMiddleware) VerifyPermission(values ...interface{}) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		log := clog.GetContextLog(c.Context())
//		if len(values)%2 != 0 {
//			log.Error().Msgf("invalid permission %v", values)
//			return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.Permission})
//		}
//
//		mapPermission := make(map[string]string)
//		for i := 0; i < len(values); i += 2 {
//			tempPem := values[i].(string)
//			tempAction := values[i+1].(string)
//			mapPermission[tempPem] = tempAction
//		}
//
//		session := jwtManager.GetSessionData(c)
//		var pass bool
//		for k, v := range mapPermission {
//			if v == common.ActionPermissionView {
//				if _, ok := session.Permission[k]; ok {
//					pass = true
//					break
//				}
//			}
//			if _, ok := session.Permission[k][v]; ok {
//				pass = true
//				break
//			}
//		}
//
//		if !pass {
//			return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.Permission})
//		}
//
//		return c.Next()
//	}
//}
//
//func (m AuthMiddleware) VerifyKYC() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		userLogin := jwtManager.GetSessionData(c)
//		if userLogin.KycLevel != common.KYCLevel1 && userLogin.KycLevel != common.KYCLevel2 {
//			if userLogin.KycStatus == common.KYCStatusPendingCode {
//				return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Verify.Kyc})
//			}
//			return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Validation.KycRequired})
//		}
//		return c.Next()
//	}
//}
//
//type VerifyPermission struct {
//	Permission map[string][]string
//}
//
//func (m AuthMiddleware) VerifyPermissionV2(verify VerifyPermission) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		session := jwtManager.GetSessionData(c)
//		var pass bool
//
//		for key, action := range verify.Permission {
//			if _, isPermission := session.Permission[key]; isPermission {
//				//if _, isAction := session.Permission[key][common.ActionPermissionView]; isAction {
//				//	pass = true
//				//	break
//				//}
//				if len(action) == 0 {
//					if _, isAction := session.Permission[key][common.ActionPermissionView]; isAction {
//						pass = true
//						break
//					}
//				}
//				for _, v := range action {
//					if _, isAction := session.Permission[key][v]; isAction {
//						pass = true
//						break
//					}
//				}
//			}
//		}
//
//		if !pass {
//			return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.Permission})
//		}
//
//		return c.Next()
//	}
//}
//
//func (m AuthMiddleware) VerifySuperAdminWhitelistIP() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		log := clog.GetContextLog(c.Context())
//		superAdminIP := c.IP()
//		if ip := c.Get("Cf-Connecting-Ip"); ip != "" {
//			superAdminIP = ip
//		}
//
//		rp := repository.NewCommonDBRepository(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)
//		whitelists, err := rp.CommonDBRepo.FindWhitelistIPSuperAdmin(c.Context())
//		if err != nil {
//			log.Error().Err(err).Msg("FindWhitelistIPSuperAdmin got err")
//			return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//		for _, whitelist := range whitelists {
//			if superAdminIP == whitelist.IP {
//				return c.Next()
//			}
//		}
//
//		log.Error().Msgf("IP %v not match ", superAdminIP)
//		return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.Permission})
//	}
//}
//
//func (m AuthMiddleware) VerifyWhitelistIP(roles ...common.UserType) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//
//		if config.Conf.State != config.StateProd {
//			return c.Next()
//		}
//
//		log := clog.GetContextLog(c.Context())
//		userLogin := jwtManager.GetSessionData(c)
//		userIP := c.IP()
//		if ip := c.Get("Cf-Connecting-Ip"); ip != "" {
//			userIP = ip
//		}
//		var err error
//
//		for _, role := range roles {
//			var whitelists []commondb.WhitelistBo
//			switch role {
//			case common.UserTypeSuperAdmin:
//				rp := repository.NewCommonDBRepository(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)
//				whitelists, err = rp.CommonDBRepo.FindWhitelistIpByRole(c.Context(), common.UserTypeSuperAdmin)
//				if err != nil {
//					log.Error().Err(err).Msg("FindWhitelistIpByRole got err")
//					return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//				}
//				for _, whitelist := range whitelists {
//					if userIP == whitelist.IP {
//						return c.Next()
//					}
//				}
//			case common.UserTypeAdmin: //For now is no requirement for use
//				rp := repository.NewCommonDBRepository(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)
//				whitelists, err = rp.CommonDBRepo.FindWhitelistIpByRole(c.Context(), common.UserTypeAdmin)
//				if err != nil {
//					log.Error().Err(err).Msg("FindWhitelistIpByRole got err")
//					return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//				}
//			case common.UserTypeSupport:
//				rp := repository.NewCommonDBRepository(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)
//				user, err := rp.CommonDBRepo.FindOneSupportByID(c.Context(), userLogin.UserID)
//				if err != nil {
//					log.Error().Err(err).Msg("FindOneSupportByID got err")
//					return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//				}
//				for _, whitelistIP := range user.WhitelistIP {
//					if userIP == whitelistIP {
//						return c.Next()
//					}
//				}
//			}
//		}
//
//		log.Error().Msgf("IP %v is not whitelist", userIP)
//		return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.Permission})
//	}
//}

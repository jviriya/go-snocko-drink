package middleware

import (
	"go-pentor-bank/internal/repository/commondb"
)

type AuthSignature struct {
	repo *commondb.SecureRepository
}

func NewAuthSignature(repo *commondb.SecureRepository) AuthSignature {
	return AuthSignature{
		repo: repo,
	}
}

//func (sv *AuthSignature) AuthenticateSignature(c *fiber.Ctx) error {
//	log := clog.GetContextLog(c.Context())
//
//	accessKey := c.Get(common.AccessKeyHeader)
//	receiveSignature := c.Get(common.SignatureHeader)
//
//	bReq := &bytes.Buffer{}
//	err := json.Compact(bReq, c.Body())
//	if err != nil {
//		log.Error().Err(err).Msg("json.Compact got err")
//		return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//	}
//	body := fmt.Sprintf("%d%s", bReq.Len(), accessKey)
//
//	openAPI, err := sv.repo.FindOneOpenAPIByAccessKey(c.Context(), accessKey)
//	if err != nil {
//		log.Error().Err(err).Msg("FindOneOpenAPIByAccessKey got err")
//		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
//	}
//
//	if openAPI.Status == common.StatusInActive {
//		return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Auth.InactiveStatus})
//	}
//
//	user, err := sv.repo.FindOneUserByUID(c.Context(), openAPI.UID)
//	if err != nil {
//		log.Error().Err(err).Msg("FindOneUserByID got err")
//		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.NotFound})
//	}
//	c.Locals(commondb.UserDataContext, user)
//
//	ip := GetIP(c)
//	mapIP := make(map[string]struct{})
//	for _, v := range openAPI.WhitelistIPs {
//		mapIP[v] = struct{}{}
//	}
//
//	if _, ok := mapIP[ip]; !ok {
//		log.Info().Msg("ip: " + ip)
//		log.Info().Msg("c.ip: " + c.IP())
//		log.Info().Msgf("c.ip: %v", c.IPs())
//		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.AccessRestricted})
//	}
//
//	expectSignature := utils.CalculateHMAC(body, openAPI.SecretKey)
//	if receiveSignature != expectSignature {
//		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.InvalidSignature})
//	}
//
//	return c.Next()
//}

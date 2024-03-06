package middleware

//func CheckMaintenance(rp *commondb.SecureRepository) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		ma, err := rp.FindOneMaintenance(c.Context())
//		if err != nil {
//			return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
//		}
//		if ma.Status == commondb.InactiveMADesc {
//			return appresponse.JSONResponse(c, http.StatusBadRequest, appresponse.IResponse{ErrorCode: config.EM.Internal.SystemMaintenance})
//		}
//
//		return c.Next()
//	}
//}

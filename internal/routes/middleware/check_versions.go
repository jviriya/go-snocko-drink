package middleware

import (
	"go-pentor-bank/internal/common"
	"strconv"
	"strings"
)

type resp struct {
	Message     string
	Description string
}

var respMap = map[string]map[string]resp{
	"maintain": {
		common.LangThai: resp{
			Message:     "ไม่สามารถใช้งานได้ในขณะนี้",
			Description: "แอปพลิเคชันไม่สามารถใช้งานได้ในขณะนี้ โปรดลองใหม่อีกครั้ง",
		},
		common.LangEnglish: resp{
			Message:     "Service unavailable.",
			Description: "Temporarily unavailable. Try again later.",
		},
	},
	"update-app": {
		common.LangThai: resp{
			Message:     "มีการอัปเดตใหม่",
			Description: "ขณะนี้เวอร์ชันใหม่พร้อมใช้งาน โปรดอัปเดตเป็นเวอร์ชัน %s",
		},
		common.LangEnglish: resp{
			Message:     "Update Available.",
			Description: "A new version is available. Please update to version %s now.",
		},
	},
	"update-os": {
		common.LangThai: resp{
			Message:     "แอปพลิเคชันรองรับตั้งแต่เวอร์ชัน %s ขึ้นไป",
			Description: "โปรดอัปเดตเวอร์ชัน %s เพื่อใช้งานแอปพลิเคชันนี้",
		},
		common.LangEnglish: resp{
			Message:     "This application requires %s or later.",
			Description: "You must update to %s in order to download and use this application.",
		},
	},
}

//func CheckVersion(rp *commondb.SecureRepository) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		//log := clog.GetContextLog(c.Context())
//
//		platform, appVersion, _ := c.Get(common.PlatformContext), c.Get(common.AppVersionContext), c.Get(common.OsVersionContext)
//		lang := common.LangEnglish
//
//		if c.Locals(common.LocaleContext) == common.LangThai {
//			lang = common.LangThai
//		}
//
//		if platform == common.WebPlatform {
//			return c.Next()
//		}
//
//		if appVersion == "" {
//			return c.Next()
//		}
//
//		//if osVersion == "" || appVersion == "" {
//		//	log.Error().Msg("osVersion is empty")
//		//	return appresponse.JSONResponse(c, http.StatusForbidden, appresponse.IResponse{})
//		//}
//
//		_, _ = rp.FindAllAppVersion(c.Context())
//		if a, f := config.Conf.AppVersions[platform]; f {
//			//if !cmpVersionCurrentGreaterThanEqual(a.OSVersion, osVersion) {
//			//	return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{
//			//		ErrorCode: config.EM.Validation.ForceUpdateOs,
//			//		VersionControl: common.AppVersionResp{
//			//			Message:     fmt.Sprintf(respMap["update-os"][lang].Message, a.OSVersion),
//			//			Description: fmt.Sprintf(respMap["update-os"][lang].Description, a.OSVersion),
//			//		},
//			//	})
//			//}
//
//			if !cmpVersionCurrentGreaterThanEqual(a.AppForceUnder, appVersion) {
//				return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{
//					ErrorCode: config.EM.Validation.ForceUpdateApplication,
//					VersionControl: common.AppVersionResp{
//						StoreLink:   a.StoreLink,
//						ForceUpdate: true,
//						Message:     respMap["update-app"][lang].Message,
//						Description: fmt.Sprintf(respMap["update-app"][lang].Description, a.AppVersion),
//					},
//				})
//			} else if !cmpVersionCurrentGreaterThanEqual(a.AppVersion, appVersion) {
//				//deviceID := c.Get(common.DeviceIdHeader)
//				//
//				//deviceDetail, err := rp.FindOneDeviceID(c.Context(), deviceID)
//				//if err != nil {
//				//	log.Error().Err(err).Msg("FindOneDeviceID got err")
//				//	return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{
//				//		ErrorCode: config.EM.Internal.InternalServerError,
//				//	})
//				//}
//
//				resp := common.AppVersionResp{
//					StoreLink:   a.StoreLink,
//					ForceUpdate: a.ForceUpdate,
//					Message:     respMap["update-app"][lang].Message,
//					Description: fmt.Sprintf(respMap["update-app"][lang].Description, a.AppVersion),
//				}
//				if a.ForceUpdate {
//					return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{
//						ErrorCode:      config.EM.Success,
//						VersionControl: resp,
//					})
//				}
//
//				c.Locals(common.AppVersionControlContext, resp)
//				return c.Next()
//			}
//		}
//
//		return c.Next()
//	}
//}

func cmpVersionCurrentGreaterThanEqual(targetVersion, currentVersion string) bool {
	currentGreaterThan := false

	if len(targetVersion) == 0 || len(currentVersion) == 0 {
		return false
	}

	tvSp := strings.Split(targetVersion, ".")
	curSp := strings.Split(currentVersion, ".")

	for index, _ := range tvSp {
		tvInt, err := strconv.Atoi(tvSp[index])
		if err != nil {
			return false
		}

		if len(curSp)-1 >= index {
			curInt, err := strconv.Atoi(curSp[index])
			if err != nil {
				return false
			}

			if tvInt < curInt {
				return true
			} else if tvInt == curInt {
				currentGreaterThan = true
			} else {
				return false
			}
		}
	}

	return currentGreaterThan
}

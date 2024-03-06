package common

import "time"

// cache
const (
	IdIPRedisKey               = "IP"
	PhoneRedisKey              = "PHONE"
	RateLimitRedisKey          = "RATE_LIMIT_"
	RateLimitExpiredAtRedisKey = "RATE_LIMIT_EXPIRED_AT_"

	RedisIdIp                              = "IP"
	LoginWebRateLimitRedisKey              = "LOGIN_WEB_RATE_LIMIT"
	LoginMobileRateLimitRedisKey           = "LOGIN_MOBILE_RATE_LIMIT"
	ValidateLoginRateLimitRedisKey         = "VALIDATE_LOGIN_RATE_LIMIT"
	RegisterRateLimitRedisKey              = "REGISTER_RATE_LIMIT"
	ValidateRegisterRateLimitRedisKey      = "VALIDATE_REGISTER_RATE_LIMIT"
	ResetPasswordRateLimitRedisKey         = "RESET_PASSWORD_RATE_LIMIT"
	ValidateResetPasswordRateLimitRedisKey = "VALIDATE_RESET_PASSWORD_RATE_LIMIT"
	SocketRoomRedisKey                     = "SOCKET_ROOM_"

	CaptchaRedisKey         = "CAPTCHA_"
	LoginCheckRedisKey      = "LOGIN_CHECK"
	CaptchaLoginWebKey      = "LOGIN_WEB"
	CaptchaLoginMobileKey   = "LOGIN_MOBILE"
	CaptchaRegisterKey      = "REGISTER"
	CaptchaResetPasswordKey = "RESET_PASSWORD"

	CacheUiTextRedisKey = "CACHE_UI_TEXT_"

	SettingConditionKey = "SETTING_CONDITION"
	SettingOpenApiKey   = "SETTING_OPEN_API"

	OtpAgeRedisExpired      = 5 * time.Minute
	OtpCoolDownRedisExpired = 1 * time.Minute
	//OtpLimitationPeriodRedisExpired = 3 * time.Minute // TODO: Must set back to 30 mins
	OtpInvalidRedisExpired   = 3 * time.Minute
	EmailOtpRedisExpired     = 5 * time.Minute
	EmailOtpRedisExpiredText = "5"
	PhoneNoOtpRedisExpired   = 5 * time.Minute
	CheckVerifyRedisExpired  = 5 * time.Minute
	CaptchaLoginRedisExpired = 1 * time.Minute

	UserOnlineRedisKey   = "USER_ONLINE_"
	UserOnlineKeyExpired = 24 * time.Hour

	CdnImgCache       = "CDN_IMG_CACHE"
	CdnImgCacheExpire = 1 * time.Hour

	CacheRoomRedisExpired = 5 * time.Minute

	AlertRedisKey = "CHECK_ALERT_"
)

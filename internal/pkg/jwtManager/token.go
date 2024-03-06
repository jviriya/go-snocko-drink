package jwtManager

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/infra/redis"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionData struct {
	UserID         primitive.ObjectID             `json:"userID,omitempty"`
	UID            string                         `json:"uid,omitempty"`
	Email          string                         `json:"email,omitempty"`
	PhoneNo        string                         `json:"phoneNo,omitempty"`
	Permission     map[string]map[string]struct{} `json:"permission,omitempty"`
	KycLevel       string                         `json:"kycLevel,omitempty"`
	KycStatus      string                         `json:"kycStatus,omitempty"`
	KycName        string                         `json:"kycName,omitempty"`
	LastSignOnDate string                         `json:"lastSignOnDate,omitempty"`
	CreatedTime    time.Time                      `json:"createdTime,omitempty"`
	Username       string                         `json:"username,omitempty"`
	NickName       string                         `json:"nickName,omitempty"`
	CountryCode    string                         `json:"countryCode,omitempty"`
	ExpTime        time.Time                      `json:"expTime,omitempty"`
	Role           common.UserType                `json:"role,omitempty"`
	UserStatus     string                         `json:"userStatus"`
	DialCode       string                         `json:"dialCode,omitempty"`
	FcmToken       string                         `json:"fcmToken"`
	Lang           string                         `json:"lang"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

type (
	JWTManager struct {
		RedisClient *redis.Redis
	}
)

func NewJWTManager(RedisClient *redis.Redis) *JWTManager {
	return &JWTManager{
		RedisClient: RedisClient,
	}
}

const (
	AuthorizationHeader = "Authorization"
	//PrefixRedis     string        = "token_"
	AtDefaultExpire time.Duration = time.Hour * 24
	RtDefaultExpire time.Duration = time.Hour * 24 * 7
)

func GetSessionData(c *fiber.Ctx) SessionData {
	if c.Locals("userLogin") == nil {
		return SessionData{}
	}
	return c.Locals("userLogin").(SessionData)
}

func GetSessionDataV2(c context.Context) SessionData {
	if c.Value("userLogin") == nil {
		return SessionData{}
	}
	return c.Value("userLogin").(SessionData)
}

func CheckSessionExist(c context.Context) (SessionData, bool) {
	if c.Value("userLogin") == nil {
		return SessionData{}, false
	}
	return c.Value("userLogin").(SessionData), true
}

func (svr *JWTManager) CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(AtDefaultExpire).Unix()
	td.AccessUUID = uuid.New().String()
	td.RtExpires = time.Now().Add(RtDefaultExpire).Unix()
	td.RefreshUUID = uuid.New().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["accessUUID"] = td.AccessUUID
	atClaims["userID"] = userID
	//atClaims["uid"] = UID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(config.Conf.ServerSetting.AccessSecret))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refreshUUID"] = td.RefreshUUID
	rtClaims["userID"] = userID
	//rtClaims["uid"] = UID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(config.Conf.ServerSetting.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (svr *JWTManager) CreateAuth(c context.Context, td *TokenDetails, data *SessionData) error {
	platform := svr.getPlatform(c)
	accessKey := fmt.Sprintf("%v_%v_access", platform, data.UserID.Hex())

	//delete same user token
	lastAccessUUID, err := svr.RedisClient.Get(c, accessKey)
	if err == nil && lastAccessUUID != "" {
		_, err := svr.RedisClient.Del(c, lastAccessUUID)
		if err != nil {
			return err
		}
	}

	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()

	data.ExpTime = at

	sessionData, err := sonic.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}

	accessUUID := td.AccessUUID
	errAccess := svr.RedisClient.Set(c, accessUUID, string(sessionData), at.Sub(now))
	if errAccess != nil {
		return errAccess
	}

	//update last token of user
	if err := svr.RedisClient.Set(c, accessKey, td.AccessUUID, at.Sub(now)); err != nil {
		return err
	}

	return nil
}

func (svr *JWTManager) ExtractToken(c context.Context, bearToken string) string {
	strArr := strings.Split(bearToken, "Bearer ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (svr *JWTManager) VerifyToken(c context.Context, bearToken string) (*jwt.Token, error) {
	tokenString := svr.ExtractToken(c, bearToken)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Conf.ServerSetting.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (svr *JWTManager) TokenValid(c context.Context, bearToken string) error {
	token, err := svr.VerifyToken(c, bearToken)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

type AccessDetails struct {
	AccessUUID string
	UserID     string
	UID        string
}

func (svr *JWTManager) ExtractTokenMetadata(c context.Context, bearToken string) (*AccessDetails, error) {
	token, err := svr.VerifyToken(c, bearToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["accessUUID"].(string)
		if !ok {
			return nil, errors.New("invalid accessUUID")
		}
		userID, ok := claims["userID"].(string)
		if !ok {
			return nil, errors.New("invalid userID")
		}
		//uid, ok := claims["uid"].(string)
		//if !ok {
		//	return nil, errors.New("invalid uid")
		//}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
			//UID:        uid,
		}, nil
	}
	return nil, err
}

func (svr *JWTManager) FetchAuth(c context.Context, accessUUID string) (SessionData, error) {
	str, err := svr.RedisClient.Get(c, accessUUID)
	if err != nil {
		return SessionData{}, err
	}
	var userLogin SessionData
	err = sonic.Unmarshal([]byte(str), &userLogin)
	if err != nil {
		return SessionData{}, err
	}
	return userLogin, nil
}

func (svr *JWTManager) DeleteAuth(c context.Context, accessUUID, userID string) error {
	platform := svr.getPlatform(c)
	_, err := svr.RedisClient.Del(c, accessUUID)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%v_%v_access", platform, userID)
	_, err = svr.RedisClient.Del(c, key)
	if err != nil {
		return err
	}
	return nil
}

func (svr *JWTManager) DeleteAllAccess(c context.Context, log *zerolog.Logger, userID string) error {
	var toDeleteUUID []string
	{
		accessKeyApp := fmt.Sprintf("%v_%v_access", common.MobilePlatform, userID)
		accessUuidApp, err := svr.RedisClient.Get(c, accessKeyApp)
		if err != nil {
			if err != goredis.Nil {
				log.Err(err).Msg("get accessKeyApp to delete got err")
			}
		} else {
			toDeleteUUID = append(toDeleteUUID, accessUuidApp)
		}

	}

	{
		accessKeyWeb := fmt.Sprintf("%v_%v_access", common.WebPlatform, userID)
		accessUuidWeb, err := svr.RedisClient.Get(c, accessKeyWeb)
		if err != nil {
			if err != goredis.Nil {
				log.Err(err).Msg("get accessUuidWeb to delete got err")
			}
		} else {
			toDeleteUUID = append(toDeleteUUID, accessUuidWeb)
		}
	}

	_, err := svr.RedisClient.Del(c, toDeleteUUID...)
	if err != nil {
		return err
	}

	return nil
}

func (svr *JWTManager) CheckLastAccessUUID(c context.Context, accessUUID, userID string) (bool, error) {
	platform := svr.getPlatform(c)
	key := fmt.Sprintf("%v_%v_access", platform, userID)
	lastAccessUUID, err := svr.RedisClient.Get(c, key)
	if err != nil {
		return false, err
	}
	return accessUUID == lastAccessUUID, nil
}

//func (svr *JWTManager) GetSessionDataByUserID(c context.Context, userID primitive.ObjectID) (string, *SessionData, error) {
//	platform := svr.getPlatform(c)
//	key := fmt.Sprintf("%v_%v_access", platform, userID.Hex())
//	lastAccessUUID, err := svr.RedisClient.Get(c, key)
//	if err != nil {
//		return "", nil, err
//	}
//
//	var data SessionData
//	err = svr.RedisClient.Parse(c, lastAccessUUID, &data)
//	if err != nil {
//		return "", nil, err
//	}
//
//	return lastAccessUUID, &data, nil
//}

type UserSessionData struct {
	LastAccessUUID string
	SessionData    SessionData
	Exist          bool
}

func (svr *JWTManager) GetAllSessionDataByUserID(c context.Context, userID primitive.ObjectID) (mobile UserSessionData, web UserSessionData, err error) {
	suffix := fmt.Sprintf("_%v_access", userID.Hex())

	var mobileData SessionData
	var exist bool
	lastAccessUUID, err := svr.RedisClient.Get(c, common.MobilePlatform+suffix)
	if err == nil {
		err = svr.RedisClient.Parse(c, lastAccessUUID, &mobileData)
		if err == nil {
			exist = true
		}
	}
	mobile = UserSessionData{
		LastAccessUUID: lastAccessUUID,
		SessionData:    mobileData,
		Exist:          exist,
	}

	var webData SessionData
	exist = false
	lastAccessUUID, err = svr.RedisClient.Get(c, common.WebPlatform+suffix)
	if err == nil {
		err = svr.RedisClient.Parse(c, lastAccessUUID, &webData)
		if err == nil {
			exist = true
		}
	}
	web = UserSessionData{
		LastAccessUUID: lastAccessUUID,
		SessionData:    webData,
		Exist:          exist,
	}

	return mobile, web, nil
}

func (svr *JWTManager) UpdateUserSessionDataByAccessUUID(c context.Context, accessUUID string, data *SessionData) error {
	sessionData, err := sonic.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}
	now := time.Now()

	errAccess := svr.RedisClient.Set(c, accessUUID, string(sessionData), data.ExpTime.Sub(now))
	if errAccess != nil {
		return errAccess
	}

	return nil
}

func (svr *JWTManager) getPlatform(c context.Context) string {
	rawPlatform := c.Value(common.PlatformContext)
	if rawPlatform != nil {
		return rawPlatform.(string)
	}

	return ""
}

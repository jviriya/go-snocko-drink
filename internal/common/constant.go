package common

const (
	LocaleContext            = "locale"
	PlatformContext          = "platform"
	RequestIDContext         = "request-id"
	AppVersionControlContext = "app-version-control"
	AppVersionContext        = "app-version"
	OsVersionContext         = "os-version"

	WebPlatform     = "web"
	MobilePlatform  = "mobile"
	IosPlatform     = "ios"
	AndroidPlatform = "android"
	BOPlatform      = "bo"

	AppAPI        = "api"
	AppBO         = "bo"
	AppOpenAPI    = "openapi"
	AppCDN        = "cdn"
	AppBatch      = "batch"
	AppSocketio   = "socket"
	AppSocketIOV2 = "socket.v2"
	AppConsumer   = "consumer"

	AccessKeyHeader = "X-Access-Key"
	SignatureHeader = "X-Signature"
	DeviceIdHeader  = "Device-Id"
)

type (
	UserType     string
	DeeplinkType string
)

const (
	DateTimeFormat1 = "02/01/2006 15:04:05"
	DateTimeFormat2 = "2006-01-02 15:04:05"
	DateTimeFormat3 = "2006-01-02 15:04:05.000"
	TimeFormat      = "15:04:05"
	TimeFormat2     = "15:04"
	DateFormat1     = "02/01/2006"
	DateFormat2     = "2006/01/02"
	DateFormat3     = "2006-01-02"
	DateFormat4     = "20060102150405"
	DateFormat5     = "060102"
	DateFormat6     = "20060102"
	DateFormat7     = "200601"
	DateFormat8     = "02-01-2006"
	DateElkFormat   = "200601"
)

const (
	RandomUsernameNumber = 8
	RandomPasswordNumber = 12
	RandomUID            = 5
	ReasonTypeOther      = "OTHER"
)

const (
	StatusActive      = "ACTIVE"
	StatusInActive    = "INACTIVE"
	StatusClosed      = "CLOSED"
	StatusAll         = "ALL"
	StatusRequired    = "REQUIRED"
	StatusNotRequired = "NOTREQUIRED"

	PhoneNoType = "PHONE"
	EmailType   = "EMAIL"

	WithdrawToUID = "UID"
)

// KYC
const (
	KYCLevel0 = "0"
	KYCLevel1 = "1"
	KYCLevel2 = "2"
)

type (
	SubscribeType string
)

// Firebase
const (
	FirebaseTypeSubscribe   SubscribeType = "SUBSCRIBE"
	FirebaseTypeUnSubscribe SubscribeType = "UNSUBSCRIBE"

	DefaultTopic          = "topic_default_%s_%s"
	DefaultTopicNotMember = "topic_default_not_member_%s_%s"
)

// Lang
const (
	LangThai    = "th"
	LangEnglish = "en"
	LangChina   = "ch"
)

// platForm
const (
	PlatformTypeWEB    = "web"
	PlatformTypeMobile = "mobile"
)

// file
const (
	ErrInvalidFileExtension       = "invalid file extension"
	ErrInvalidFileMIME            = "invalid file MIME"
	ErrExtensionNotRequired       = "file extension is not required"
	ErrFileSizeLimitation         = "file size exceeds the limitation"
	ErrInvalidExtensionOnFilename = "invalid file extension compared with original file"
)

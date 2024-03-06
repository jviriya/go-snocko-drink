package config

import (
	"fmt"
	"go-pentor-bank/internal/common"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/wawafc/go-utils/money"

	"github.com/spf13/viper"
)

const (
	validatorRequired         = "required"
	validatorMin              = "min"
	validatorMax              = "max"
	validatorNotEqual         = "ne"
	validatorLen              = "len"
	validatorOneOf            = "oneof"
	validatorGreaterThan      = "gt"
	validatorGreaterThanEqual = "gte"
	validatorLessThan         = "lt"
	validatorLessThanEqual    = "lte"

	validateEmail       = "email"
	validateNumeric     = "numeric"
	validateDate        = "date"
	validateThai        = "thaialpha"
	validateAlphaNum    = "alphanum"
	validateNationality = "thaination"
	validateMoney       = "money"
	validateCurrency    = "currency"
	validateMobile      = "mobile"
)

var (
	mobileRegex  = regexp.MustCompile("^([0]{1}|(66))[0-9]{9,9}$")
	englishRegex = regexp.MustCompile("^[a-z,A-Z]*$")
	Validate     *validator.Validate
)

type LocaleDescription struct {
	Locale string
	TH     string
	EN     string
	CH     string
}

type ErrorCode struct {
	Code    string            `json:"code"`
	Message LocaleDescription `json:"description"`
	TraceID string            `json:"traceID"`
}

func (ec *ErrorCode) IsSuccess() bool {
	if ec.Code == "0000" {
		return true
	}
	return false
}

func (ec *ErrorCode) WithLocale(c *fiber.Ctx) ErrorCode {
	locale := c.Locals(common.LocaleContext)
	if locale == nil || locale == "" {
		ec.Message.Locale = "en"
	} else {
		ec.Message.Locale = locale.(string)
	}
	return *ec
}

func (ec *ErrorCode) WithFormat(a ...interface{}) ErrorCode {
	ec.Message.TH = fmt.Sprintf(ec.Message.TH, a...)
	ec.Message.EN = fmt.Sprintf(ec.Message.EN, a...)

	return *ec
}

func (ld LocaleDescription) MarshalJSON() ([]byte, error) {
	switch strings.ToLower(ld.Locale) {
	case "th":
		return sonic.Marshal(ld.TH)
	case "ch":
		return sonic.Marshal(ld.CH)
	default:
		return sonic.Marshal(ld.EN)
	}
}

var EM ErrorMessage

type ErrorMessage struct {
	vn         *viper.Viper
	ConfigPath string

	Success  ErrorCode
	Internal struct {
		General               ErrorCode
		BadRequest            ErrorCode
		InternalServerError   ErrorCode
		DatabaseError         ErrorCode
		Timeout               ErrorCode
		RequestLimit          ErrorCode
		InvalidReclaimBalance ErrorCode
		PermissionDenied      ErrorCode
		SystemMaintenance     ErrorCode
	}
	Auth struct {
		NotFound                    ErrorCode
		IncorrectPassword           ErrorCode
		Permission                  ErrorCode
		AttemptedUnauthorizedAccess ErrorCode
		IncorrectPhone              ErrorCode
		InactiveAccount             ErrorCode
		InvalidPlatform             ErrorCode
		InvalidSignature            ErrorCode
		AccessRestricted            ErrorCode
		IncorrectUsernameOrPassword ErrorCode
	}
	Validation struct {
		InvalidPasswordCondition     ErrorCode
		InvalidEmployeeCodeCondition ErrorCode
		InvalidFirstNameCondition    ErrorCode
		InvalidLastNameCondition     ErrorCode
		ValidationFailed             ErrorCode
		InvalidEmailPattern          ErrorCode
		InvalidPhonePattern          ErrorCode
		InsufficientBalance          ErrorCode
		FirebaseCannotRegisTopic     ErrorCode
		//InsufficientTokenQuantity             ErrorCode
		//OrderTimeout                          ErrorCode
		//InvalidCurrency                       ErrorCode
		//InvalidTimeFormat                     ErrorCode
		//OtpExceedLimitation                   ErrorCode
		//OtpWaitingPeriod                      ErrorCode
		//TransferOwnAccount                    ErrorCode
		//InsufficientFundsWallet               ErrorCode
		//ChangeNicknameOverThan3               ErrorCode
		//NotAvailableInYourCountry             ErrorCode
		//ConditionNotMatch                 ErrorCode
		//InvalidTransferStatus                 ErrorCode
		//AdCannotHaveMoreThanLimit             ErrorCode
		//AdCannotActiveMoreThanLimit           ErrorCode
		//InvalidToken                          ErrorCode
		//InvalidQrData                         ErrorCode
		//InvalidTransferExpTime                ErrorCode
		//InvalidTransferQuantity               ErrorCode
		//InvalidTransferStatusSuccess          ErrorCode
		//InvalidTransferStatusExpired          ErrorCode
		//InvalidTransferStatusCanceled         ErrorCode
		//InvalidOrderTokenAmount               ErrorCode
		//InvalidDialCode                       ErrorCode
		//EndDateMustBeAfterStartDate           ErrorCode
		//InvalidDateRange                      ErrorCode
		//InvalidOtpType                        ErrorCode
		//RequiredOtpNewType                    ErrorCode
		//InvalidNotificationType               ErrorCode
		//AdOffline                             ErrorCode
		//InvalidSubscriptionId                 ErrorCode
		//InvalidSubscriptionPlan               ErrorCode
		//InvalidStartDateEndDate               ErrorCode
		//InvalidSubscriptionToken              ErrorCode
		//AlreadyReview                         ErrorCode
		//ServerTransferReject                  ErrorCode
		//DestinationTransferCanNotReceiveToken ErrorCode
		//AdQuantityLessThanLimit               ErrorCode
		//OrderIsAlreadyExpired                 ErrorCode
		//ForceUpdateApplication                ErrorCode
		//AlertExceedLimit                      ErrorCode
		//WarningUpdateApplication              ErrorCode
	}
}

type ValidateStructResponse struct {
	Field         string      `json:"field"`
	Tag           string      `json:"tag,omitempty"`
	ExpectedValue interface{} `json:"expectedValue"`
	ReceivedValue interface{} `json:"receivedValue"`
	Message       string      `json:"message,omitempty"`
	MapField      string      `json:"mapField"`
}

func (em *ErrorMessage) Init(errorPath string) error {
	em.ConfigPath = errorPath

	vn := viper.New()
	vn.AddConfigPath(errorPath)
	vn.SetConfigName("error")

	if err := vn.ReadInConfig(); err != nil {
		return err
	}

	em.vn = vn

	em.mapping("", reflect.ValueOf(em).Elem())

	return nil
}

func (em *ErrorMessage) mapping(name string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fi := v.Field(i)
		if fi.Kind() != reflect.Struct {
			continue
		}

		fn := Underscore(v.Type().Field(i).Name)
		if name != "" {
			fn = fmt.Sprint(name, ".", fn)
		}

		if fi.Type().Name() == "ErrorCode" {
			fi.Set(reflect.ValueOf(em.ErrorCode(fn)))
			continue
		}
		em.mapping(fn, fi)
	}
}

func (em *ErrorMessage) ErrorCode(name string) ErrorCode {
	rtn := ErrorCode{
		Code: em.vn.GetString(fmt.Sprintf("%s.code", name)),
		Message: LocaleDescription{
			TH: em.vn.GetString(fmt.Sprintf("%s.th", name)),
			EN: em.vn.GetString(fmt.Sprintf("%s.en", name)),
			CH: em.vn.GetString(fmt.Sprintf("%s.ch", name)),
		},
	}
	return rtn
}

func Underscore(str string) string {
	runes := []rune(str)
	var out []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 && (unicode.IsUpper(runes[i]) || unicode.IsNumber(runes[i])) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

func ValidateStruct(v interface{}) ([]*ValidateStructResponse, error) {
	var errors []*ValidateStructResponse
	err := Validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidateStructResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.ExpectedValue = err.Param()
			element.ReceivedValue = err.Value()
			errors = append(errors, &element)
		}
	}
	return errors, err
}

func InitDefaultValidators() error {
	if Validate == nil {
		Validate = validator.New(validator.WithRequiredStructEnabled())
		if err := Validate.RegisterValidation("date", dateValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("thaialpha", thaiAlphaValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("englishalpha", engAlphaValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("mobile", mobileValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("pin", ValidatePin); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("status", isValidStatus); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("isLimit50", isLimit50); err != nil {
			return err
		}
		Validate.RegisterCustomTypeFunc(moneyCustomTypeFunc, money.Money{})
		Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	return nil
}

func moneyCustomTypeFunc(field reflect.Value) interface{} {
	if value, ok := field.Interface().(money.Money); ok {
		return value.Float64()
	}
	return nil
}

func dateValidator(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006/01/02", fl.Field().String()); err != nil {
		return false
	}
	return true
}

func thaiAlphaValidator(fl validator.FieldLevel) bool {
	s := []rune(fl.Field().String())

	for _, r := range s {
		if !unicode.Is(unicode.Thai, r) {
			return false
		}
	}
	return true
}

func engAlphaValidator(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return englishRegex.MatchString(v)
}

func mobileValidator(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return mobileRegex.MatchString(v)
}

func IsValidPasswordCondition(s string) bool {
	var letters, number bool
	if utf8.RuneCountInString(s) < 8 {
		return false
	}
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLetter(c):
			letters = true
		}
	}
	return letters && number
}

func isLimit50(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if utf8.RuneCountInString(v) > 50 {
		return false
	}
	return true
}

func isValidStatus(fl validator.FieldLevel) bool {

	v := fl.Field().String()
	switch strings.ToUpper(v) {
	case common.StatusActive, common.StatusInActive:
		return true
	default:
		return false
	}
}

func ValidateTypeUtilsMoney(amount money.Money) bool {
	if amount.GreaterThan(money.NewMoneyFromFloat(0)) {
		return true
	}
	return false
}

func ValidatePin(fl validator.FieldLevel) bool {

	v := fl.Field().String()
	if len(v) != 4 {
		return false
	}
	for _, c := range v {
		if !unicode.IsNumber(c) {
			return false
		}
	}

	return true
}

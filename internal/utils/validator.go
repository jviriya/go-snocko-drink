package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/wawafc/go-utils/money"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	ValidatorRequired         = "required"
	ValidatorMin              = "min"
	ValidatorMax              = "max"
	ValidatorNotEqual         = "ne"
	ValidatorLen              = "len"
	ValidatorOneOf            = "oneof"
	ValidatorGreaterThan      = "gt"
	ValidatorGreaterThanEqual = "gte"
	ValidatorLessThan         = "lt"
	ValidatorLessThanEqual    = "lte"

	ValidateEmail       = "email"
	ValidateNumeric     = "numeric"
	ValidateDate        = "date"
	ValidateThai        = "thaialpha"
	ValidateAlphaNum    = "alphanum"
	ValidateNationality = "thaination"
	ValidateMoney       = "money"
	ValidateCurrency    = "currency"
	ValidateMobile      = "mobile"
)

var (
	engCharAndNumber    = regexp.MustCompile("^[a-zA-Z0-9]+$")
	mobileRegex         = regexp.MustCompile("^([0]{1}|(66))[0-9]{9,9}$")
	englishRegex        = regexp.MustCompile("^[a-z,A-Z]*$")
	Validate            *validator.Validate
	nicknameCharAllowed = []rune{'_', '-', '@', '.', '(', ')', '*'}
)

func InitDefaultValidators() error {
	if Validate == nil {
		Validate = validator.New(validator.WithRequiredStructEnabled())
		if err := Validate.RegisterValidation("date", dateValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("thaialpha", thaiAlphaValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("englishalpha", thaiAlphaValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("mobile", mobileValidator); err != nil {
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
	var upper, lower, number bool
	if utf8.RuneCountInString(s) < 8 {
		return false
	}
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLower(c):
			upper = true
		case unicode.IsUpper(c):
			lower = true
		}
	}
	return upper && lower && number
}

func IsValidNameCondition(s string) bool {
	for _, c := range s {
		if !unicode.IsNumber(c) {
			continue
		} else {
			return false
		}
	}
	return true
}

func IsValidUsernameCondition(s string) bool {
	for _, c := range s {
		switch {
		case unicode.IsNumber(c), unicode.IsLower(c), unicode.IsUpper(c):
			continue
		default:
			return false
		}
	}
	return true
}

func IsValidEmployeeCodeCondition(s string) bool {
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			continue
		default:
			return false
		}
	}
	return true
}
func IsValidMasterDataString(s string, master []string) bool {
	for _, d := range master {
		if d == s {
			return true
		}
	}
	return false
}

func IsValidNumberCondition(s string) bool {
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			continue
		default:
			return false
		}
	}
	return true
}

func IsValidYearOfServiceCondition(s string) bool {
	var symbol, number bool

	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsSymbol(c):
			symbol = true
		}
	}
	return symbol && number
}

func IsEmailValid(email string) bool {
	// Define the regular expression pattern for a valid email address
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Use the regexp package to compile the pattern into a regular expression object
	regex := regexp.MustCompile(pattern)

	// Use the MatchString method of the regular expression object to test the email address
	return regex.MatchString(email)
}

func IsMobileValid(mobile string) bool {
	pattern := `^\d+$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(mobile)
}

func IsNicknameCharacterValid(nickname string) bool {
	for _, char := range nickname {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && !unicode.IsSpace(char) && !unicode.IsMark(char) {
			ok := slices.Contains(nicknameCharAllowed, char)
			if !ok {
				return false
			}
		}
	}
	return true
}

func IsOnlyCharAndNumber(data string) bool {
	return engCharAndNumber.MatchString(data)
}

func TransferExpiredTime(transExpTime *time.Time) bool {

	if transExpTime != nil {
		if time.Now().After(*transExpTime) {
			return true
		}
	}
	return false
}

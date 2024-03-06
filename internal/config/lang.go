package config

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/common"
	"reflect"
	"strings"
)

const (
	LangEN = "en"
	LangTH = "th"
	LangCN = "cn"
	LangLA = "la"
	LangID = "id"
	LangVN = "vn"
	LangMY = "my"
	LangMM = "mm"
	LangTW = "tw"
)

type Lang struct {
	En string `json:"en,omitempty" bson:"en,omitempty" abbr:"en" fullName:"English"`
	Th string `json:"th,omitempty" bson:"th,omitempty" abbr:"th" fullName:"Thai"`
	Cn string `json:"cn,omitempty" bson:"cn,omitempty" abbr:"cn" fullName:"Chinese"`
	La string `json:"la,omitempty" bson:"la,omitempty" abbr:"la" fullName:"Lao"`
	Id string `json:"id,omitempty" bson:"id,omitempty" abbr:"id" fullName:"Indonesian"`
	Vn string `json:"vn,omitempty" bson:"vn,omitempty" abbr:"vn" fullName:"Vietnamese"`
	My string `json:"my,omitempty" bson:"my,omitempty" abbr:"my" fullName:"Malay"`
	Mm string `json:"mm,omitempty" bson:"mm,omitempty" abbr:"mm" fullName:"Burmese (Myanmar)"`
	Kh string `json:"kh,omitempty" bson:"kh,omitempty" abbr:"kh" fullName:"Khmer"`
	Tw string `json:"tw,omitempty" bson:"tw,omitempty" abbr:"tw" fullName:"Mandarin Chinese"`
}

// 'THB','th' = ไทย
// 'USA','en' = อังกฤษ
// 'CNY','cn' = จีน
// 'IDR','id' = อินโดนีเซีย
// 'VND','vn' = เวียดนาม
// 'LAK','la' = ลาว
// 'MYR','my' = มาเลเซีย
// 'MMK','mm' = พม่า
// 'KHR','kh' = กัมพูชา

func ValidateLanguage(lang string) bool {
	switch lang {
	case LangEN, LangTH, LangCN, LangLA, LangID, LangVN, LangMY, LangMM, LangTW:
		return true
	}
	return false
}

func (l *Lang) SetByLang(lang string, value string) {
	langData := reflect.ValueOf(l).Elem()
	for i := 0; i < langData.NumField(); i++ {
		if langData.Type().Field(i).Tag.Get("abbr") == lang {
			f := langData.Field(i)
			if f.CanSet() {
				f.SetString(value)
			}
			break
		}
	}
}

func (l *Lang) SetValueOfLang(lang Lang) {
	langData := reflect.ValueOf(&lang).Elem()
	langMap := make(map[string]string)
	for i := 0; i < langData.NumField(); i++ {
		langMap[langData.Type().Field(i).Name] = langData.Field(i).String()
	}
	lData := reflect.ValueOf(l).Elem()
	for i := 0; i < lData.NumField(); i++ {
		f := lData.FieldByName(lData.Type().Field(i).Name)
		if f.CanSet() {
			if lData.Field(i).String() == "" {
				f.SetString(langMap[lData.Type().Field(i).Name])
			}
		}
	}
}

func (l *Lang) LangToString(lang string) string {
	data := l.En
	langData := reflect.ValueOf(*l)
	for i := 0; i < langData.NumField(); i++ {
		if langData.Type().Field(i).Tag.Get("abbr") == lang && langData.Field(i).String() != "" {
			data = langData.Field(i).String()
			break
		}
	}
	return data
}

func (l *Lang) ByLocaleContext(c context.Context) string {
	if c.Value(common.LocaleContext) == nil {
		return ""
	}
	return l.LangToString(c.Value(common.LocaleContext).(string))
}

func GetLang(c *fiber.Ctx) string {
	if c.Locals(common.LocaleContext) == nil {
		return "en"
	}
	return c.Locals(common.LocaleContext).(string)
}

func GamatronLanguage(lang string) string {
	switch strings.ToLower(lang) {
	case "cn":
		return "zh"
	case "vn":
		return "vi"
	default:
		return "en"
	}
}

func GetLanguageList() ([]string, []string) {
	ref := reflect.ValueOf(Lang{})

	var langs []string
	var names []string
	for i := 0; i < ref.NumField(); i++ {
		tags := ref.Type().Field(i).Tag.Get("abbr")
		fullName := ref.Type().Field(i).Tag.Get("fullName")
		tag := strings.Split(tags, ",")
		langs = append(langs, tag[0])
		names = append(names, fullName)
	}

	return langs, names
}

func getLanguageValue(src Lang) map[string]string {
	ref := reflect.ValueOf(src)

	data := make(map[string]string)
	for i := 0; i < ref.NumField(); i++ {
		key := ref.Type().Field(i).Name
		val := ref.Field(i).String()
		data[key] = val
	}
	return data
}

func SetLanguageStruct(src Lang, dest interface{}) {
	srcData := getLanguageValue(src)

	ref := reflect.ValueOf(dest).Elem()
	for key, val := range srcData {
		ref.FieldByName(key).SetString(val)
	}
}

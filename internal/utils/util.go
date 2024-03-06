package utils

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"go-pentor-bank/internal/config"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"

	b64 "encoding/base64"
	//smSession "go-pentor-bank/internal/infra/sm"
)

var Stage = config.Conf.State

func PrintJson(data interface{}) {
	if config.Conf.State != config.StateProd {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		if err := enc.Encode(data); err != nil {
			panic(err)
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func GenMD5(parameter string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(parameter)))
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func InquiryGetSkipLimit(input map[string]interface{}) (int64, int64) {
	var page int64 = 1
	if val, ok := input["page"]; ok {
		page = int64(val.(float64))
	}
	var size int64 = 25
	if val, ok := input["size"]; ok {
		size = int64(val.(float64))
	}
	skip := (page - 1) * size
	return skip, size
}

func InquiryGetSort(input map[string]interface{}, defaultSort string, defaultBy string, specialCase map[string]string) bson.D {
	sortBy := defaultBy
	if by, validateType := input["by"].(string); validateType && (by == "desc" || by == "asc") {
		sortBy = input["by"].(string)
	}
	by := 1
	if sortBy == "desc" {
		by = -1
	}

	sortresult := bson.D{}
	if sortValue, validType1 := input["sort"].(string); validType1 {
		if nameInDbS, validType2 := specialCase[sortValue]; validType2 {
			nameInDbL := strings.Split(nameInDbS, "|")
			for _, v := range nameInDbL {
				sortresult = append(sortresult, bson.E{v, by})
			}
		} else if sortValue != "" {
			sortresult = append(sortresult, bson.E{sortValue, by})
		}
	}
	if len(sortresult) == 0 {
		sortresult = append(sortresult, bson.E{defaultSort, by})
	}
	return sortresult
}

func trimMapStringInterface(data interface{}) interface{} {
	if values, valid := data.([]interface{}); valid {
		for i := range values {
			data.([]interface{})[i] = trimMapStringInterface(values[i])
		}
	} else if values, valid := data.(map[string]interface{}); valid {
		for k, v := range values {
			data.(map[string]interface{})[k] = trimMapStringInterface(v)
		}
	} else if value, valid := data.(string); valid {
		data = strings.TrimSpace(value)
	}
	return data
}

func validateCitizenID(cid string) error {
	if len(cid) == 0 {
		return errors.New("please enter a valid citizen ID/certified document ID")
	}

	if len(cid) != 13 {
		return errors.New("please enter a valid citizen ID/certified document ID")
	}

	//check sum
	sum := 0
	for i := 0; i < 12; i++ {
		number, err := strconv.Atoi(string(cid[i]))
		if err != nil {
			return err
		}
		sum += number * (13 - i)
	}

	c, err := strconv.Atoi(string(cid[12]))
	if err != nil {
		return err
	}

	if (11-sum%11)%10 != c {
		return errors.New("please enter a valid citizen ID/certified document ID")
	}

	return nil
}

func GetIntParam(input map[string]interface{}, name string) (int, error) {
	data, found := input[name]
	if !found {
		return 0, errors.New("not found")
	}
	value, valid := data.(int)
	if !valid {
		return 0, errors.New("invalid type")
	}
	return value, nil
}

func GetInt64Param(input map[string]interface{}, name string) (int64, error) {
	data, found := input[name]
	if !found {
		return 0, errors.New("not found")
	}
	value, valid := data.(int64)
	if !valid {
		return 0, errors.New("invalid type")
	}
	return value, nil
}

func GetFloat64Param(input map[string]interface{}, name string) (float64, error) {
	data, found := input[name]
	if !found {
		return 0, errors.New("not found")
	}
	value, valid := data.(float64)
	if !valid {
		return 0, errors.New("invalid type")
	}
	return value, nil
}

func GetStringParam(input map[string]interface{}, name string) (string, error) {
	data, found := input[name]
	if !found {
		return "", errors.New("not found")
	}
	value, valid := data.(string)
	if !valid {
		return "", errors.New("invalid type")
	}
	return value, nil
}

func GetBoolParam(input map[string]interface{}, name string) (bool, error) {
	data, found := input[name]
	if !found {
		return false, errors.New("not found")
	}
	value, valid := data.(bool)
	if !valid {
		return false, errors.New("invalid type")
	}
	return value, nil
}

func ReqToTime(startDate, startTime, endDate, endTime, format, formatTime string, timeZone *time.Location) (time.Time, time.Time, error) {
	if startTime == "" {
		startTime = "00:00:00"
	}
	start, err := time.ParseInLocation(formatTime, fmt.Sprintf(format, startDate, startTime), timeZone)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	var end time.Time
	if endTime == "" {
		end, err = time.ParseInLocation(formatTime, fmt.Sprintf(format, endDate, "00:00:00"), timeZone)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		end = end.AddDate(0, 0, 1)
	} else {
		end, err = time.ParseInLocation(formatTime, fmt.Sprintf(format, endDate, endTime), timeZone)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		end = end.Add(time.Second)
	}
	return start, end, nil
}

func EndTimeOfDay(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.In(loc).Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, loc).Add(-1)
}

func InterfaceToString(data interface{}) string {
	dataStr := ""
	if v, valid := data.(string); valid {
		dataStr = v
	} else if v, valid := data.(float64); valid {
		dataStr = fmt.Sprintf("%0.f", v)
	} else if v, valid := data.(int); valid {
		dataStr = fmt.Sprintf("%d", v)
	} else if v, valid := data.(int64); valid {
		dataStr = fmt.Sprintf("%d", v)
	} else if v, valid := data.(int32); valid {
		dataStr = fmt.Sprintf("%d", v)
	} else if v, valid := data.(float32); valid {
		dataStr = fmt.Sprintf("%0.f", v)
	}
	return dataStr
}

func InterfaceToBool(data interface{}) bool {
	dataBool := false
	if rtv, valid := data.(string); valid {
		switch rtv {
		case "TRUE":
			dataBool = true
		case "FALSE":
			dataBool = false
		}
	} else if rtv, valid := data.(bool); valid {
		dataBool = rtv
	}
	return dataBool
}

func MakeTimestampMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func MakeMillisecondToTimestamp(millisecond int64) time.Time {
	return time.Unix(0, millisecond*int64(time.Millisecond))
}

func MapStrToListKey(m map[string]interface{}) []interface{} {
	keys := []interface{}{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func EscapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

type InquiryPattern struct {
	FieldNameFront string
	FieldNameDB    string
	Type           string
	Operation      string
}

func InquiryGenBsonM(fillter map[string]interface{}, patterns []InquiryPattern) bson.M {
	selector := bson.M{}
	for _, pattern := range patterns {
		switch pattern.Operation {
		case "eq":
			if value, found := fillter[pattern.FieldNameFront]; found && isset(value) {
				switch pattern.Type {
				case "objectID":
					if valueStr, valid := value.(string); valid {
						ObjID, _ := primitive.ObjectIDFromHex(valueStr)
						selector[pattern.FieldNameDB] = ObjID
					} else {
						selector[pattern.FieldNameDB] = value
					}
				default:
					selector[pattern.FieldNameDB] = value
				}
			}
		case "like":
			if value, found := fillter[pattern.FieldNameFront]; found && isset(value) {
				if valueStr, valid := value.(string); valid {
					selector[pattern.FieldNameDB] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + regexp.QuoteMeta(valueStr) + ".*", Options: "i"}}
				}
			}

		case "range":
			switch pattern.Type {
			case "date":
				dateParam := strings.Split(pattern.FieldNameFront, "|")
				dateFillter := bson.M{}
				if value, found := fillter[dateParam[0]]; found {
					if valueStr, valid := value.(string); valid {
						fromDate, err := time.ParseInLocation("2006-01-02", valueStr, config.TimeZone.Bangkok)
						if err == nil {
							dateFillter["$gte"] = fromDate
						}
					}
				}
				if len(dateParam) > 1 {
					if value, found := fillter[dateParam[1]]; found {
						if valueStr, valid := value.(string); valid {
							toDate, err := time.ParseInLocation("2006-01-02", valueStr, config.TimeZone.Bangkok)
							if err == nil {
								toDate = toDate.AddDate(0, 0, 1)
								dateFillter["$lt"] = toDate
							}
						}
					}
				}
				if len(dateFillter) > 0 {
					selector[pattern.FieldNameDB] = dateFillter
				}
			default:
				dateParam := strings.Split(pattern.FieldNameFront, "|")
				numFillter := bson.M{}
				if value, found := fillter[dateParam[0]]; found {
					numFillter["$gte"] = value
				}
				if len(dateParam) > 1 {
					if value, found := fillter[dateParam[1]]; found {
						numFillter["$lte"] = value
					}
				}
				if len(numFillter) > 0 {
					selector[pattern.FieldNameDB] = numFillter
				}
			}
		}
	}
	return selector
}

func isset(input interface{}) bool {
	return input != nil && input != ""
}

func FormatDate(loc *time.Location, t *time.Time, format string) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.In(loc).Format(format)
}

//func getKeyEncryptionKey() []byte {
//	return []byte(smSession.Conn.GLT2022key)
//}
//
//func GLT2022Encrypt(plaintext string) string {
//	switch strings.ToLower(string(Stage)) {
//	case "prod", "production":
//		key := getKeyEncryptionKey()
//		return glt32EncryptSP(key, plaintext)
//	default:
//		return plaintext
//	}
//}
//
//func GLT2022Decrypt(cipherText string) string {
//	switch strings.ToLower(string(Stage)) {
//	case "prod", "production":
//		key := getKeyEncryptionKey()
//		return glt32DecryptSP(key, cipherText)
//	default:
//		return cipherText
//	}
//}

func glt32EncryptSP(key []byte, plaintext string) string {
	b64Text := strings.Replace(base64Encode(plaintext), "=", "*", -1)
	var prepareAesText string

	if len(b64Text) >= aes.BlockSize {
		prepareAesText = b64Text[:aes.BlockSize]
		cipherText := encryptAES(key, prepareAesText)
		return fmt.Sprintf("%v%v", cipherText, b64Text[aes.BlockSize:])
	} else {
		padding := specialPaddingString(b64Text)
		cipherText := encryptAES(key, padding)
		return cipherText
	}
}

func glt32DecryptSP(key []byte, cipherText string) string {
	//i := strings.Index(cipherText, ">")
	if len(cipherText) > 32 {
		aesText := cipherText[:32]
		rest := cipherText[32:]

		b64Shard := decryptAES(key, aesText)
		ori, _ := base64Decode(strings.Replace(b64Shard+rest, "*", "=", -1))

		return ori
	} else {
		peep := decryptAES(key, cipherText)
		oriN := strings.SplitN(peep, "-", 2)
		ori, _ := base64Decode(strings.Replace(oriN[0], "*", "=", -1))

		return ori
	}
}

func specialPaddingString(ori string) string {
	if len(ori) > 16 {
		return ori
	} else {
		rand.Seed(time.Now().UnixNano())
		//pad := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=-")
		pad := []rune("+") //fix to single rune
		min := 0
		max := len(pad)

		ori = fmt.Sprintf("%v-", ori)
		for len(ori) < 16 {
			ori = fmt.Sprintf("%v%v", ori, string(pad[rand.Intn(max-min)+min]))
		}

		return ori
	}
}

func base64Encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, bool) {
	data, err := b64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", true
	}
	return string(data), false
}

func encryptAES(key []byte, plaintext string) string {

	c, err := aes.NewCipher(key)
	CheckError(err)

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func decryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	CheckError(err)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s
}

func CheckError(err error) {
	if err != nil {
		log.Err(err).Msg("crypto go err")
	}
}

func GetCounter(currentPage, pageLimit int64) int64 {
	return pageLimit * (currentPage - 1)
}

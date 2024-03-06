package utils

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName() (string, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	_, more := frames.Next()
	if !more {
		return "", ""
	}
	_, more = frames.Next()
	if !more {
		return "", ""
	}
	frame, _ := frames.Next()
	str := strings.Split(frame.Function, "/")
	strPNT := strings.Split(str[len(str)-1], ".")
	packageName := strPNT[0]
	funcName := strPNT[len(strPNT)-1]
	return packageName, funcName
}

func ToLowerMemberUsername(username, appID string) string {
	return strings.ToLower(fmt.Sprintf("%v@%v", username, appID))
}

func GetStructTag(f reflect.StructField, tagName string) string {
	return f.Tag.Get(tagName)
}

func GetOnlyPhoneNumber(phoneNo string) string {
	if len(phoneNo) > 1 && strings.HasPrefix(phoneNo, "0") {
		return phoneNo[1:]
	}
	return phoneNo
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsIPv4(ipv4 string) bool {
	ip := net.ParseIP(ipv4)

	return ip != nil && ip.To4() != nil
}

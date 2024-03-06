package utils

import (
	"fmt"
	"reflect"
	"time"
)

func StructToQueryParams(data interface{}) string {
	var queryParam string
	e := reflect.ValueOf(data)
	for i := 0; i < e.NumField(); i++ {
		tagName := e.Type().Field(i).Tag.Get("json")
		var value string
		switch e.Field(i).Interface().(type) {
		case time.Time:
			value = fmt.Sprint(e.Field(i).Interface().(time.Time).UnixNano())
		default:
			value = fmt.Sprint(e.Field(i).Interface())
		}
		queryParam += tagName + "=" + value + "&"
	}
	return queryParam[:len(queryParam)-1]
}

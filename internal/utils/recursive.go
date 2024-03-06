package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
)

func Walk(v reflect.Value) {
	v = v.Elem()
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			Walk(v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			Walk(v.MapIndex(k))
		}
	case reflect.Struct:
		switch v.Interface().(type) {
		case primitive.Decimal128:
			v.Field(v.NumField()).SetFloat(100555.0)
			fmt.Println(v.Interface())
			fmt.Printf("Visiting %v type %v\n", v, v.Kind())
		default:
		}
	default:
	}
}

func WalkMap(t interface{}) interface{} {
	switch t.(type) {
	case []interface{}:
		data := t.([]interface{})
		for i, v := range data {
			data[i] = WalkMap(v)
		}
		t = data
	case primitive.A:
		data := t.(primitive.A)
		for i, v := range data {
			data[i] = WalkMap(v)
		}
		t = data
	case map[string]interface{}:
		data := t.(map[string]interface{})
		for k, v := range data {
			data[k] = WalkMap(v)
		}
		t = data
	case primitive.Decimal128:
		data := t.(primitive.Decimal128)
		if s, err := strconv.ParseFloat(data.String(), 64); err == nil {
			return s
		}
		return t
	default:
		return t
	}
	return t
}

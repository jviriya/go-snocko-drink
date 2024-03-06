package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/wawafc/go-utils/money"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func readEnvironmentConfig(v interface{}) {
	reflectType := reflect.TypeOf(v)
	reflectValue := reflect.ValueOf(v)

	if reflectType.Kind() == reflect.Ptr {
		//fmt.Printf("Argument is a pointer, dereferencing.\n")
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	for i := 0; i < reflectType.NumField(); i++ {
		types := reflectType.Field(i).Type
		field := reflectValue.Field(i)

		var envName string
		var envVal string
		if temp := reflectValue.Type().Field(i).Tag.Get("overrideEnv"); temp != "" {
			envVal = os.Getenv(temp)
			envName = temp
		}

		switch reflectValue.Field(i).Kind() {
		case reflect.Slice:
			if envName != "" {
				decodeSlice(reflectValue.Field(i), envVal)
			}
		case reflect.Array:
			//panic("not support array")
		case reflect.Bool:
			if envName != "" {
				newVal, err := strconv.ParseBool(envVal)
				if err != nil {
					panic(err)
				}
				field.SetBool(newVal)
			}
		case reflect.String:
			if envName != "" {
				field.SetString(envVal)
			}
		case reflect.Int32, reflect.Int64, reflect.Int:
			if envName != "" {
				newVal, err := strconv.ParseInt(envVal, 10, 64)
				if err != nil {
					panic(fmt.Errorf("envName: %s got err %v", envName, err))
				}
				field.SetInt(newVal)
			}
		case reflect.Float32, reflect.Float64:
			if envName != "" {
				newVal, err := strconv.ParseFloat(envVal, 64)
				if err != nil {
					panic(err)
				}
				field.SetFloat(newVal)
			}
		case reflect.Struct:
			valueInf := field.Interface()
			switch valueInf.(type) {
			case viper.Viper:
				// ignore
				fmt.Println("ignore Viper")
			case money.Money:
				if envName != "" {
					s, _ := money.NewMoneyFromString(os.Getenv(envName))
					rs := reflect.ValueOf(&s).Elem()
					newVal := reflect.NewAt(types, unsafe.Pointer(rs.UnsafeAddr())).Elem()
					field.Set(newVal)
				}
			default:
				readEnvironmentConfig(field.Addr().Interface())
			}
		}
	}
}

func decodeSlice(val reflect.Value, envVal string) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	envSlice := strings.Split(envVal, ",")
	val.SetZero()
	newReflect := reflect.Indirect(reflect.New(val.Type()))
	for i, v := range envSlice {
		for newReflect.Len() <= i {
			newReflect = reflect.Append(newReflect, reflect.Zero(val.Type().Elem()))
		}
		switch val.Type().Elem().Kind() {
		case reflect.Bool:
			newVal, err := strconv.ParseBool(v)
			if err != nil {
				panic(err)
			}
			newReflect.Index(i).SetBool(newVal)
		case reflect.String:
			newReflect.Index(i).SetString(strings.TrimSpace(v))
		case reflect.Int32, reflect.Int64, reflect.Int:
			newVal, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				panic(err)
			}
			newReflect.Index(i).SetInt(newVal)
		case reflect.Float32, reflect.Float64:
			newVal, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}
			newReflect.Index(i).SetFloat(newVal)
		}
	}
	val.Set(newReflect)
}

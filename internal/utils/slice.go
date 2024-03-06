package utils

import "strings"

func RemoveDuplicate[T string | ~int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func TrimSliceSpace(sliceList []string) []string {
	for i := 0; i < len(sliceList); i++ {
		sliceList[i] = strings.TrimSpace(sliceList[i])
	}

	return sliceList
}

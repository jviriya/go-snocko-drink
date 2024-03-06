package utils

import (
	"fmt"
)

func ToAppIDString(keyDate string, runningNo int64) string {
	return fmt.Sprintf("%s%09d", keyDate, runningNo)
}

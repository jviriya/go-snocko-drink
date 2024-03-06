package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenOTP(digit int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	otp := ""
	for i := 0; i < digit; i++ {
		otp += strconv.Itoa(rand.Intn(9))
	}
	return otp
}

func GenKeyVerifyOTP(phone string) string {
	return fmt.Sprintf("lead_votp_ref_%s", phone)
}
func GenKeyRecoveryOTP(phone string) string {
	return fmt.Sprintf("pos_rotp_ref_%s", phone)
}

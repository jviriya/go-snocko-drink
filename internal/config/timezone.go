package config

import "time"

type TimeLocation struct {
	Bangkok   *time.Location
	Singapore *time.Location
}

var TimeZone TimeLocation

func LoadTimeLocation() error {

	bangkokTZ, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	TimeZone.Bangkok = bangkokTZ

	singaporeTZ, err := time.LoadLocation("Singapore")
	if err != nil {
		return err
	}
	TimeZone.Singapore = singaporeTZ

	return nil
}

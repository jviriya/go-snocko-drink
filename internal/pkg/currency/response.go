package currency

import (
	"github.com/wawafc/go-utils/money"
	"time"
)

type RatesResp1 struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates struct {
		USD money.Money `json:"USD"`
		CNY money.Money `json:"CNY"`
		IDR money.Money `json:"IDR"`
		KHR money.Money `json:"KHR"`
		LAK money.Money `json:"LAK"`
		MMK money.Money `json:"MMK"`
		MYR money.Money `json:"MYR"`
		THB money.Money `json:"THB"`
		VND money.Money `json:"VND"`
		AED money.Money `json:"AED"`
		AUD money.Money `json:"AUD"`
		BHD money.Money `json:"BHD"`
		BND money.Money `json:"BND"`
		CAD money.Money `json:"CAD"`
		CHF money.Money `json:"CHF"`
		DKK money.Money `json:"DKK"`
		EUR money.Money `json:"EUR"`
		GBP money.Money `json:"GBP"`
		HKD money.Money `json:"HKD"`
		JPY money.Money `json:"JPY"`
		INR money.Money `json:"INR"`
		KRW money.Money `json:"KRW"`
		NOK money.Money `json:"NOK"`
		NZD money.Money `json:"NZD"`
		PHP money.Money `json:"PHP"`
		QAR money.Money `json:"QAR"`
		RUB money.Money `json:"RUB"`
		SAR money.Money `json:"SAR"`
		SEK money.Money `json:"SEK"`
		SGD money.Money `json:"SGD"`
		TWD money.Money `json:"TWD"`
		ZAR money.Money `json:"ZAR"`
	} `json:"rates"`
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
	Unit      string `json:"unit"`
}

type RatesResp struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data struct {
		USD RateDetail `json:"USD"`
		CNY RateDetail `json:"CNY"`
		IDR RateDetail `json:"IDR"`
		KHR RateDetail `json:"KHR"`
		LAK RateDetail `json:"LAK"`
		MMK RateDetail `json:"MMK"`
		MYR RateDetail `json:"MYR"`
		THB RateDetail `json:"THB"`
		VND RateDetail `json:"VND"`
		AED RateDetail `json:"AED"`
		AUD RateDetail `json:"AUD"`
		BHD RateDetail `json:"BHD"`
		BND RateDetail `json:"BND"`
		CAD RateDetail `json:"CAD"`
		CHF RateDetail `json:"CHF"`
		DKK RateDetail `json:"DKK"`
		EUR RateDetail `json:"EUR"`
		GBP RateDetail `json:"GBP"`
		HKD RateDetail `json:"HKD"`
		JPY RateDetail `json:"JPY"`
		INR RateDetail `json:"INR"`
		KRW RateDetail `json:"KRW"`
		NOK RateDetail `json:"NOK"`
		NZD RateDetail `json:"NZD"`
		PHP RateDetail `json:"PHP"`
		QAR RateDetail `json:"QAR"`
		RUB RateDetail `json:"RUB"`
		SAR RateDetail `json:"SAR"`
		SEK RateDetail `json:"SEK"`
		SGD RateDetail `json:"SGD"`
		TWD RateDetail `json:"TWD"`
		ZAR RateDetail `json:"ZAR"`
	} `json:"data"`
}

type RateDetail struct {
	Code  string      `json:"code"`
	Value money.Money `json:"value"`
}

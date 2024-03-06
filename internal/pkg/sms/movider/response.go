package movider

type Response struct {
	RemainingBalance float64 `json:"remaining_balance"`
	TotalSms         int     `json:"total_sms"`
	PhoneNumberList  []struct {
		Number    string  `json:"number"`
		MessageID string  `json:"message_id"`
		TotalSms  int     `json:"total_sms"`
		Price     float64 `json:"price"`
	} `json:"phone_number_list"`
	BadPhoneNumberList []any `json:"bad_phone_number_list"`
	Error              struct {
		Code        int    `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"error"`
}

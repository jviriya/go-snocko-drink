package vonage

type SendSMSResp struct {
	MessageCount string `json:"message-count"`
	Messages     []struct {
		To               string `json:"to"`
		MessageId        string `json:"message-id"`
		Status           string `json:"status"`
		RemainingBalance string `json:"remaining-balance"`
		MessagePrice     string `json:"message-price"`
		Network          string `json:"network"`
		ClientRef        string `json:"client-ref"`
		AccountRef       string `json:"account-ref"`
	} `json:"messages"`
}

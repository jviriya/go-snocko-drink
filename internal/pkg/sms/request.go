package sms

type SendSMSReq struct {
	Tel  string `json:"tel"`
	Text string `json:"text"`
}

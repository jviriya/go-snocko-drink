package deecommerce

type SendSMSReq struct {
	AccountID string `json:"accountId"`
	SecretKey string `json:"secretKey"`
	Type      string `json:"type"`
	To        string `json:"to"`
	Sender    string `json:"sender"`
	Msg       string `json:"msg"`
}

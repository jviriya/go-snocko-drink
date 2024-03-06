package deecommerce

type SendSMSResp struct {
	Error        string `json:"error"`
	Msg          string `json:"msg"`
	DeliveryCode string `json:"delivery_code"`
	QuotaBalance int    `json:"quota_balance"`
}

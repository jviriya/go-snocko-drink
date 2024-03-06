package curl

func GenerateHeaderUrlencoded() map[string]string {
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	return header
}

func GenerateHeaderXML() map[string]string {
	header := map[string]string{
		"Content-Type": "application/xml",
		//"request-id":     requestid,
		//"request-app-id": appid,
	}
	return header
}

func GenerateHeaderJson() map[string]string {
	header := map[string]string{
		"Content-Type": "application/json",
		//"request-id":     requestid,
		//"request-app-id": appid,
	}
	return header
}

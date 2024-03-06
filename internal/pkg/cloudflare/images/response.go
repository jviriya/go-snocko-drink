package images

import "time"

type UploadImageResp struct {
	Result struct {
		ID                string    `json:"id"`
		Filename          string    `json:"filename"`
		Uploaded          time.Time `json:"uploaded"`
		RequireSignedURLs bool      `json:"requireSignedURLs"`
		Variants          []string  `json:"variants"`
	} `json:"result"`
	ResultInfo any     `json:"result_info"`
	Success    bool    `json:"success"`
	Errors     []Error `json:"errors"`
	Messages   []Error `json:"messages"`
}

type ListImages struct {
	Result struct {
		Images []struct {
			ID                string    `json:"id"`
			Filename          string    `json:"filename"`
			Uploaded          time.Time `json:"uploaded"`
			RequireSignedURLs bool      `json:"requireSignedURLs"`
			Variants          []string  `json:"variants"`
		} `json:"images"`
	} `json:"result"`
	Success  bool  `json:"success"`
	Errors   []any `json:"errors"`
	Messages []any `json:"messages"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

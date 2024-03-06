package appresponse

import (
	"github.com/bytedance/sonic"
	"go-pentor-bank/internal/config"
	"time"
)

type QueryOptional struct {
	CurrentPage  int64 `json:"currentPage,omitempty"`
	TotalPages   int64 `json:"totalPages"`
	PageLimit    int64 `json:"pageLimit"`
	TotalRecords int64 `json:"totalRecords"`
	//Counter      int64 `json:"counter"`
	//Limit        int    `json:"limit,omitempty"`
	//Offset       int    `json:"offset,omitempty"`
	//Sort         []Sort `json:"sort,omitempty"`
}

type Sort struct {
	By    string `json:"by,omitempty"`
	Order int    `json:"order,omitempty"`
}

type IResponse struct {
	config.ErrorCode
	Error           error          `json:"error,omitempty"`
	ValidationError interface{}    `json:"validationError,omitempty"`
	ResponseTime    string         `json:"responseTime"`
	Data            interface{}    `json:"data,omitempty"`
	VersionControl  interface{}    `json:"versionControl"`
	QueryOptional   *QueryOptional `json:"paging,omitempty"`
}

func (i IResponse) MarshalJSON() ([]byte, error) {
	type iResponse IResponse
	resp := &struct {
		iResponse
		ErrorResp string `json:"error,omitempty"`
	}{
		iResponse: (iResponse)(i),
	}
	if i.Error != nil {
		resp.ErrorResp = i.Error.Error()
	}
	resp.ResponseTime = time.Now().Format(time.RFC3339Nano)
	return sonic.Marshal(resp)
}

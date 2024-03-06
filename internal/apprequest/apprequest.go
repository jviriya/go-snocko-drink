package apprequest

import (
	"github.com/wawafc/go-utils/money"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	PagingDefaultType QueryOtpType = ""
	NextPageType      QueryOtpType = "NEXT_PAGE"
)

type (
	QueryOtpType string
)

type QueryOptional struct {
	Type      QueryOtpType `json:"type"`
	PageNo    int64        `json:"pageNo"`
	PageLimit int64        `json:"pageLimit"`
	CurrentID string       `json:"currentID"`
	NextPage  NextPage     `json:"nextPage"`
}

type NextPage struct {
	IDAfter       *primitive.ObjectID `json:"idAfter"`
	DateTimeAfter *time.Time          `json:"dateTimeAfter"`
	PriceAfter    *money.Money        `json:"priceAfter"`
}

func (queryOpt *QueryOptional) GetSkipLimit() (int64, int64) {
	return (queryOpt.PageNo - 1) * queryOpt.PageLimit, queryOpt.PageLimit
}

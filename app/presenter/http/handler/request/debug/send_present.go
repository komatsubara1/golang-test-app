package request

import "time"

type DebugSendPresentRequest struct {
	Title       string     `json:"title" binding:"required"`
	Params      string     `json:"params" binding:"required"`
	ContentType uint       `json:"content_type" binding:"required"`
	ContentId   uint64     `json:"content_id" binding:"required"`
	Amount      uint64     `json:"amount" binding:"required"`
	ArriveDate  *time.Time `json:"arrive_date" binding:"required"`
	LimitDate   *time.Time `json:"limit_date" binding:"required"`
}

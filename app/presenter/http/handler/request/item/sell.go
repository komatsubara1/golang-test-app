package request

import (
	"app/domain/value/master"
)

type ItemSellRequest struct {
	ItemId master.ItemId `json:"item_id" binding:"required"`
	Count  uint64        `json:"count" binding:"required,min=1"`
}

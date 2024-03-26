package request

import (
	"app/domain/value/master"
)

type ItemGetRequest struct {
	ItemId master.ItemId `json:"item_id" binding:"required"`
}

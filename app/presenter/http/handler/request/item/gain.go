package request

type ItemGainRequest struct {
	ItemId uint64 `json:"item_id" binding:"required,min=1"`
	Count  uint64 `json:"count" binding:"required,min=1"`
}

package response

import (
	"app/domain/entity/user"
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewItemGainResponse(message string, code error.ErrorCode, item *user.UserItem) ItemGainResponse {
	return ItemGainResponse{message, code, item}
}

type ItemGainResponse struct {
	Message string          `binding:"required,min=1"`
	Code    error.ErrorCode `binding:"required"`
	Item    *user.UserItem  `binding:"required"`
}

func (r ItemGainResponse) ToJson() gin.H {
	return gin.H{
		"item":    r.ToJsonItem(),
		"message": r.Message,
		"code":    r.Code,
	}
}

func (r ItemGainResponse) ToJsonItem() gin.H {
	return gin.H{
		"item_id":  r.Item.ItemId.Value(),
		"quantity": r.Item.Quantity,
	}
}

package response

import (
	"app/domain/entity/user"
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewItemGetAllResponse(message string, code error.ErrorCode, items *user.UserItems) ItemGetAllResponse {
	return ItemGetAllResponse{message, code, items}
}

type ItemGetAllResponse struct {
	Message string          `binding:"required,min=1"`
	Code    error.ErrorCode `binding:"required"`
	Items   *user.UserItems `binding:"required"`
}

func (r ItemGetAllResponse) ToJson() gin.H {
	return gin.H{
		"items":   r.ToJsonItems(),
		"message": r.Message,
		"code":    r.Code,
	}
}

func (r ItemGetAllResponse) ToJsonItems() []gin.H {
	var a []gin.H

	for _, v := range *r.Items {
		a = append(a, gin.H{
			"item_id":  v.ItemId.Value(),
			"quantity": v.Quantity,
		})
	}

	return a
}

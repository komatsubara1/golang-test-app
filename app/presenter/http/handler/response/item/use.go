package response

import (
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"

	user_entity "app/domain/entity/user"
)

func NewItemUseResponse(message string, code error.ErrorCode, user *user_entity.User, item *user_entity.UserItem) ItemUseResponse {
	return ItemUseResponse{message, code, user, item}
}

type ItemUseResponse struct {
	Message string                `binding:"required,min=1"`
	Code    error.ErrorCode       `binding:"required"`
	User    *user_entity.User     `binding:"required"`
	Item    *user_entity.UserItem `binding:"required"`
}

func (r ItemUseResponse) ToJson() gin.H {
	return gin.H{
		"user":    r.ToJsonUser(),
		"item":    r.ToJsonItem(),
		"message": r.Message,
		"code":    r.Code,
	}
}

func (r ItemUseResponse) ToJsonUser() gin.H {
	return gin.H{
		"id":                        r.User.ID.Value(),
		"name":                      r.User.Name,
		"stamina":                   r.User.Stamina,
		"stamina_latest_updated_at": r.User.StaminaLatestUpdatedAt,
		"coin":                      r.User.Coin,
		"latest_logged_in_at":       r.User.LatestLoggedInAt,
	}
}

// TODO: 持たせるのEntity？
func (r ItemUseResponse) ToJsonItem() gin.H {
	return gin.H{
		"item_id":  r.Item.ItemId.Value(),
		"quantity": r.Item.Quantity,
	}
}

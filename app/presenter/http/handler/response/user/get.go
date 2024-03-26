package response

import (
	user_entity "app/domain/entity/user"
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewUserGetResponse(message string, code error.ErrorCode, user *user_entity.User) UserGetResponse {
	return UserGetResponse{message, code, user}
}

type UserGetResponse struct {
	Message string            `binding:"required,min=1"`
	Code    error.ErrorCode   `binding:"required"`
	User    *user_entity.User `binding:"required"`
}

func (r UserGetResponse) ToJson() gin.H {
	return gin.H{
		"user":    r.ToJsonUser(),
		"message": r.Message,
		"code":    r.Code,
	}
}

func (r UserGetResponse) ToJsonUser() gin.H {
	return gin.H{
		"id":                        r.User.ID.Value(),
		"name":                      r.User.Name,
		"stamina":                   r.User.Stamina,
		"stamina_latest_updated_at": r.User.StaminaLatestUpdatedAt,
		"coin":                      r.User.Coin,
		"latest_logged_in_at":       r.User.LatestLoggedInAt,
	}
}

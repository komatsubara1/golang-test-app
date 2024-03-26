package response

import (
	user_entity "app/domain/entity/user"
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewUserCreateResponse(message string, code error.ErrorCode, user *user_entity.User) UserCreateResponse {
	return UserCreateResponse{message, code, user}
}

type UserCreateResponse struct {
	Message string            `binding:"required,min=1"`
	Code    error.ErrorCode   `binding:"required"`
	User    *user_entity.User `binding:"required"`
}

func (r UserCreateResponse) ToJson() gin.H {
	return gin.H{
		"user":    r.ToJsonUser(),
		"message": r.Message,
		"code":    r.Code,
	}
}

func (r UserCreateResponse) ToJsonUser() gin.H {
	return gin.H{
		"id":                        r.User.ID.Value(),
		"name":                      r.User.Name,
		"stamina":                   r.User.Stamina,
		"stamina_latest_updated_at": r.User.StaminaLatestUpdatedAt,
		"coin":                      r.User.Coin,
		"latest_logged_in_at":       r.User.LatestLoggedInAt,
	}
}

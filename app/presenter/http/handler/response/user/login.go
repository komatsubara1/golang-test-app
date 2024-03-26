package response

import (
	userentity "app/domain/entity/user"
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewUserLoginResponse(message string, code error.ErrorCode, user *userentity.User) UserLoginResponse {
	return UserLoginResponse{message, code, user}
}

type UserLoginResponse struct {
	Message string           `binding:"required,min=1"`
	Code    error.ErrorCode  `binding:"required"`
	User    *userentity.User `binding:"required"`
}

func (r UserLoginResponse) ToJson() gin.H {
	return gin.H{
		"user":    r.ToJsonUser(),
		"message": r.Message,
		"code":    r.Code,
	}
}

func (r UserLoginResponse) ToJsonUser() gin.H {
	return gin.H{
		"id":                        r.User.ID.Value(),
		"name":                      r.User.Name,
		"stamina":                   r.User.Stamina,
		"stamina_latest_updated_at": r.User.StaminaLatestUpdatedAt,
		"coin":                      r.User.Coin,
		"latest_logged_in_at":       r.User.LatestLoggedInAt,
	}
}

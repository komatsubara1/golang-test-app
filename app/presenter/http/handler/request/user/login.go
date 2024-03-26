package request

import "app/domain/value/user"

type UserLoginRequest struct {
	UserId user.UserId `json:"user_id" binding:"required"`
}

package response

import (
	userEntity "app/domain/entity/user"
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewDebugSendPresentResponse(
	message string,
	code error.ErrorCode,
	present userEntity.UserPresent,
) DebugSendPresentResponse {
	return DebugSendPresentResponse{Message: message, Code: code, UserPresent: present}
}

type DebugSendPresentResponse struct {
	Message     string                 `binding:"required,min=1"`
	Code        error.ErrorCode        `binding:"required"`
	UserPresent userEntity.UserPresent `binding:"required"`
}

func (r DebugSendPresentResponse) ToJson() gin.H {
	return gin.H{
		"message":      r.Message,
		"code":         r.Code,
		"user_present": r.UserPresent,
	}
}

package response

import (
	"app/domain/enum/error"

	"github.com/gin-gonic/gin"
)

func NewDebugSetTimeShiftResponse(message string, code error.ErrorCode) DebugSetTimeShiftResponse {
	return DebugSetTimeShiftResponse{message, code}
}

type DebugSetTimeShiftResponse struct {
	Message string          `binding:"required,min=1"`
	Code    error.ErrorCode `binding:"required"`
}

func (r DebugSetTimeShiftResponse) ToJson() gin.H {
	return gin.H{
		"message": r.Message,
		"code":    r.Code,
	}
}

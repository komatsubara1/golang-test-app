package response

import "github.com/gin-gonic/gin"

type CommonResponse struct {
	Message string `binding:"required,min=1"`
	Code    int64  `binding:"required"`
}

func (r CommonResponse) ToJson() gin.H {
	return gin.H{
		"message": r.Message,
		"code":    r.Code,
	}
}

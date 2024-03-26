package response

import (
	"app/domain/entity/user"
	"app/domain/enum/error"
	"app/presenter/http/handler/response"
)

func NewItemGetResponse(message string, code error.ErrorCode, item *user.UserItem) ItemGetResponse {
	return ItemGetResponse{C: response.CommonResponse{Message: message, Code: int64(code)}, Item: item}
}

type ItemGetResponse struct {
	C    response.CommonResponse `json:"common"`
	Item *user.UserItem          `json:"item" binding:"required"`
}

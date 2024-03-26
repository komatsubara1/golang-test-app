package request

type UserCreateRequest struct {
	UserName string `json:"user_name" binding:"required"`
}

package service

import (
	userEntity "app/domain/entity/user"
	"app/domain/value/user"
	"time"
)

type UserAuthService struct{}

func (r UserAuthService) Build(userId user.UserId, token string, now time.Time) *userEntity.UserAuth {
	var ua = &userEntity.UserAuth{}
	ua.UserId = userId
	ua.Token = token
	ua.ExpiredAt = now.AddDate(0, 0, 1)
	return ua
}

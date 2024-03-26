package service

import (
	"app/define"
	userEntity "app/domain/entity/user"
	"app/domain/value/user"
	"time"
)

type UserService struct{}

func (r UserService) Build(userId user.UserId, name string, now time.Time) *userEntity.User {
	var u = &userEntity.User{}
	u.ID = userId
	u.Name = name
	u.Stamina = define.InitialStamina
	u.StaminaLatestUpdatedAt = now
	u.Coin = define.InitialCoin
	u.LatestLoggedInAt = now
	return u
}

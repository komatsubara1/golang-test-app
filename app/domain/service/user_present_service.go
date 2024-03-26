package service

import (
	userEntity "app/domain/entity/user"
	"app/domain/value/user"
	"github.com/google/uuid"
	"time"
)

type UserPresentService struct{}

func (r UserPresentService) Build(
	userId user.UserId,
	title string,
	params string,
	contentType uint,
	contentId uint64,
	amount uint64,
	limitDate *time.Time,
) *userEntity.UserPresent {
	up := &userEntity.UserPresent{}
	up.ID = user.NewPresentId(uuid.New())
	up.UserId = userId
	up.Title = title
	up.Params = params
	up.ContentType = contentType
	up.ContentId = contentId
	up.Amount = amount
	up.ArriveDate = nil
	up.LimitDate = limitDate
	up.ReceivedAt = nil
	return up
}

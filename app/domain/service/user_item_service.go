package service

import (
	"app/define"
	masterEntity "app/domain/entity/master"
	userEntity "app/domain/entity/user"
	"app/domain/value/master"
	userValue "app/domain/value/user"
	"github.com/ahmetalpbalkan/go-linq"
	"time"
)

type UserItemService struct{}

func (r UserItemService) Build(userId userValue.UserId, itemId master.ItemId, count uint64) *userEntity.UserItem {
	var ui = &userEntity.UserItem{}
	ui.UserId = userId
	ui.ItemId = itemId
	ui.Quantity = count
	return ui
}

func (r UserItemService) FilterIsInStock(list *userEntity.UserItems) *userEntity.UserItems {
	ret := &userEntity.UserItems{}
	linq.From(list).Where(func(item interface{}) bool {
		return r.IsInStock(item.(*userEntity.UserItem))
	}).ToSlice(&ret)

	return ret
}

func (r UserItemService) IsInStock(userItem *userEntity.UserItem) bool {
	return userItem.Quantity > 0
}

func (r UserItemService) IsExceedMaxCount(userItem *userEntity.UserItem, item *masterEntity.ItemMaster) bool {
	return userItem.Quantity > item.MaxCount
}

func (r UserItemService) RecoverStamina(
	user *userEntity.User, item *masterEntity.ItemMaster, now time.Time,
) *userEntity.User {
	autoRecoverStamina, fractionTime := r.calcAutoRecoverStamina(user.StaminaLatestUpdatedAt)
	addedStamina := user.Stamina + autoRecoverStamina + item.EffectValue
	if addedStamina < define.StaminaMax {
		user.Stamina = addedStamina
	} else {
		user.Stamina = define.StaminaMax
	}
	user.StaminaLatestUpdatedAt = now.Add(fractionTime)

	return user
}

// TODO: implements
func (r UserItemService) calcAutoRecoverStamina(staminaLatestUpdatedAt time.Time) (uint64, time.Duration) {
	return 0, time.Duration(-1)
}

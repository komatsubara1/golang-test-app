// Code generated by cli/generator/entity/gen.go; DO NOT EDIT.

package user

import (
	"app/domain/value/user"
	"time"
)

type UserPresents []UserPresent

type UserPresent struct {
	ID          user.PresentId `json:"id,string" gorm:"column:id;type:varchar(255)"`
	UserId      user.UserId    `json:"user_id,string" gorm:"column:user_id;type:varchar(255)"`
	Title       string         `json:"title" gorm:"column:title"`
	Params      string         `json:"params" gorm:"column:params"`
	ContentType uint           `json:"content_type" gorm:"column:content_type"`
	ContentId   uint64         `json:"content_id" gorm:"column:content_id"`
	Amount      uint64         `json:"amount" gorm:"column:amount"`
	ArriveDate  *time.Time     `json:"arrive_date" gorm:"column:arrive_date"`
	LimitDate   *time.Time     `json:"limit_date" gorm:"column:limit_date"`
	ReceivedAt  *time.Time     `json:"received_at" gorm:"column:received_at"`
	CreatedAt   time.Time      `json:"-" gorm:"-;autoCreateTime"`
	UpdatedAt   time.Time      `json:"-" gorm:"-;autoUpdateTime"`
}

func (e *UserPresent) TableName() string {
	return "user_present"
}

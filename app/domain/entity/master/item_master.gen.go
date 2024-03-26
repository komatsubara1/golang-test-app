// Code generated by cli/generator/entity/gen.go; DO NOT EDIT.

package master

import (
	"app/domain/value/master"
	"time"
)

type ItemMasters []ItemMaster

type ItemMaster struct {
	ID           master.ItemId     `json:"id,int" gorm:"column:id;type:int"`
	Name         string            `json:"name" gorm:"column:name"`
	Type         uint              `json:"type" gorm:"column:type"`
	SellCoin     uint64            `json:"sell_coin" gorm:"column:sell_coin"`
	EffectType   uint              `json:"effect_type" gorm:"column:effect_type"`
	EffectValue  uint64            `json:"effect_value" gorm:"column:effect_value"`
	ScheduleId   master.ScheduleId `json:"schedule_id,int" gorm:"column:schedule_id;type:int"`
	MaxCount     uint64            `json:"max_count" gorm:"column:max_count"`
	MaxViewCount uint64            `json:"max_view_count" gorm:"column:max_view_count"`
	CreatedAt    time.Time         `json:"-" gorm:"-;autoCreateTime"`
	UpdatedAt    time.Time         `json:"-" gorm:"-;autoUpdateTime"`
}

func (e *ItemMaster) TableName() string {
	return "item_master"
}
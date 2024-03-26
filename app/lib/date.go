package lib

import (
	"app/domain/value/master"
	"time"
)

type Include int

const (
	IncludeNone Include = 0
	IncludeStart
	IncludeEnd
	IncludeBoth
)

// Between 期間内
func Between(start, end, now time.Time, include Include) bool {
	return start.After(now) && end.Before(now)
}

// BetweenSchedule
// TODO: スケジュール期間内
func BetweenSchedule(scheduleId master.ScheduleId, now time.Time, include Include) bool {
	return false
}

// IsOverEndOfDay
// TODO: 日跨ぎ
func IsOverEndOfDay(start, end, dateLine time.Time) bool {
	return false
}

// IsOverEndOfWeek
// TODO: 週跨ぎ
func IsOverEndOfWeek(start, end, dateLine time.Time) bool {
	return false
}

// IsOverEndOfMonth
// TODO: 月跨ぎ
func IsOverEndOfMonth(start, end, dateLine time.Time) bool {
	return false
}

// GetElapsedDay
// TODO: 2点間の経過日数を取得
func GetElapsedDay(start, end, dateLine time.Time) bool {
	return false
}

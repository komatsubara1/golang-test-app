// Code generated by cli/generator/vo/gen.go; DO NOT EDIT.

package master

import (
	"app/domain/value"
	"encoding/json"
	"fmt"
)

type ScheduleId struct {
	value.ValueObject[int64]
}

func NewScheduleId(v int64) ScheduleId {
	return ScheduleId{value.NewValueObject[int64](v)}
}

func (v *ScheduleId) Scan(value interface{}) error {
	switch vt := value.(type) {
	case int64:
		*v = NewScheduleId(vt)

	default:
		return fmt.Errorf("invalid type. type=%s", vt)
	}
	return nil
}

func (v *ScheduleId) UnmarshalJSON(data []byte) error {
	var t any
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	switch t.(type) {
	case float64:
		v.Scan(int64(t.(float64)))
		return nil
	}

	v.Scan(t)
	return nil
}
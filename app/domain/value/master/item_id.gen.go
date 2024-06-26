// Code generated by cli/generator/vo/gen.go; DO NOT EDIT.

package master

import (
	"app/domain/value"
	"encoding/json"
	"fmt"
)

type ItemId struct {
	value.ValueObject[uint64]
}

func NewItemId(v uint64) ItemId {
	return ItemId{value.NewValueObject[uint64](v)}
}

func (v *ItemId) Scan(value interface{}) error {
	switch vt := value.(type) {
	case uint64:
		*v = NewItemId(vt)
	case int64:
		*v = NewItemId(uint64(vt))
	default:
		return fmt.Errorf("invalid type. type=%s", vt)
	}
	return nil
}

func (v *ItemId) UnmarshalJSON(data []byte) error {
	var t any
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	switch t.(type) {
	case float64:
		v.Scan(uint64(t.(float64)))
		return nil
	}

	v.Scan(t)
	return nil
}

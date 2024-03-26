package lib

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Sum[T Number](a, b T) T {
	return a + b
}

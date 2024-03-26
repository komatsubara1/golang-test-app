package lib

import (
	"app/context"
)

const RedisKey = "time_shift"

func FindTimeShift(gctx *context.GameContext) int64 {
	userId := gctx.UserId
	gctx.Ucctx.Connect()
	res := HGet[int64](gctx.Ucctx.Conn, RedisKey, userId)
	if res == nil {
		return 0
	}
	return *res
}

func SetTimeShift(gctx *context.GameContext, shiftTime int64) error {
	userId := gctx.UserId

	gctx.Ucctx.Connect()
	return HSet[int64](gctx.Ucctx.Conn, RedisKey, userId, &shiftTime)
}

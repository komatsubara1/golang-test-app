package middleware

import (
	"app/context"
	"app/lib"
	"github.com/gin-gonic/gin"
	"time"
)

const RedisKey = "time_shift"

func TimeShiftMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gctx := ctx.MustGet("GameContext").(*context.GameContext)
		timeShift := lib.FindTimeShift(gctx)

		gctx.UtcNow = gctx.UtcNow.Add(time.Duration(timeShift))
	}
}

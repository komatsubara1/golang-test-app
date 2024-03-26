package middleware

import (
	"app/context"
	"github.com/gin-gonic/gin"
	"time"
)

// SetUpMiddleware リクエスト内で必要なものをcontextに詰める
func SetUpMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("UtcNow", time.Now().UTC())

		userDbContext := context.NewUserDbContext()
		userCacheContext := context.NewUserCacheContext()
		masterDbContext := context.NewMasterDbContext()
		masterCacheContext := context.NewMasterCacheContext()
		gameContext := context.NewGameContext(
			userDbContext,
			userCacheContext,
			masterDbContext,
			masterCacheContext,
			time.Now().UTC(),
		)
		ctx.Set("GameContext", gameContext)
	}
}

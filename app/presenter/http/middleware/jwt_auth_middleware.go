package middleware

import (
	"app/context"
	"app/infrastructure/repository/user"
	"app/lib"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// JwtAuthenticationMiddleware 認証
func JwtAuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader(os.Getenv("TOKEN_KEY"))

		secretKey := os.Getenv("TOKEN_SECRET")

		userId, _, err := lib.ParseToken(tokenString, secretKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "failed parse token.",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		gctx := ctx.MustGet("GameContext").(*context.GameContext)
		udctx := gctx.Udctx
		udctx.Connect()
		repos := user.NewUserAuthRepository()
		userAuth, err := repos.FindByUserId(ctx, userId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "failed find user_auth.",
				"userId":  userId.Value(),
			})
			ctx.Abort()
			return
		}

		if userAuth == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "user_auth not found.",
				"userId":  userId,
			})
			ctx.Abort()
			return
		}

		if !(userAuth.Token == tokenString) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message":      "invalid token.",
				"userId":       userAuth.UserId,
				"token":        userAuth.Token,
				"requestToken": tokenString,
			})
			ctx.Abort()
			return
		}

		if !userAuth.ExpiredAt.Before(ctx.MustGet("UtcNow").(time.Time)) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message":      "invalid token.",
				"userId":       userAuth.UserId,
				"token":        userAuth.Token,
				"requestToken": tokenString,
			})
			ctx.Abort()
			return
		}

		gctx.UserId = &userId
		udctx.SetUserId(userId)

		ctx.Next()
	}
}

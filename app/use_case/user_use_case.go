package use_case

import (
	"app/context"
	userEntity "app/domain/entity/user"
	userRepository "app/domain/repository/user"
	"app/domain/service"
	userValue "app/domain/value/user"
	"app/lib"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserUseCase struct {
	ur  userRepository.UserRepository
	uar userRepository.UserAuthRepository
}

func NewUserUseCase(ur userRepository.UserRepository, uar userRepository.UserAuthRepository) *UserUseCase {
	return &UserUseCase{ur: ur, uar: uar}
}

// Get ユーザー情報取得
func (u *UserUseCase) Get(ctx *gin.Context, userId userValue.UserId) (*userEntity.User, error) {
	user, err := u.ur.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create ユーザー登録
func (u *UserUseCase) Create(ctx *gin.Context, name string) (*userEntity.User, *userEntity.UserAuth, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}

	userId := userValue.NewUserId(uid)
	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	ctx.Set("user_id", userId)

	tokenSecret := os.Getenv("TOKEN_SECRET")
	tokenLifetime, err := strconv.Atoi(os.Getenv("TOKEN_LIFETIME"))
	if err != nil {
		return nil, nil, err
	}

	user := service.UserService{}.Build(userId, name, ctx.MustGet("UtcNow").(time.Time))

	token, err := lib.GenerateToken(user.ID, tokenSecret, ctx.MustGet("UtcNow").(time.Time), tokenLifetime)
	if err != nil {
		return nil, nil, err
	}

	userAuth := service.UserAuthService{}.Build(user.ID, token, ctx.MustGet("UtcNow").(time.Time))

	err = udctx.TransactionScope(func() error {
		if err := u.ur.Save(ctx, *user); err != nil {
			return err
		}

		if err := u.uar.Save(ctx, *userAuth); err != nil {
			return err
		}

		return nil
	})

	return user, userAuth, err
}

// Login ログイン
func (u *UserUseCase) Login(ctx *gin.Context, userId userValue.UserId) (*userEntity.User, *userEntity.UserAuth, error) {
	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx

	tokenSecret := os.Getenv("TOKEN_SECRET")
	tokenLifetime, err := strconv.Atoi(os.Getenv("TOKEN_LIFETIME"))
	if err != nil {
		return nil, nil, err
	}

	token, err := lib.GenerateToken(userId, tokenSecret, ctx.MustGet("UtcNow").(time.Time), tokenLifetime)
	if err != nil {
		return nil, nil, err
	}

	var user *userEntity.User = nil
	var userAuth *userEntity.UserAuth = nil
	err = udctx.TransactionScope(func() error {
		user, err := u.ur.FindByID(ctx, userId)
		if err != nil {
			return err
		}

		if user == nil {
			return err
		}

		userAuth, err := u.uar.FindByUserId(ctx, userId)
		if err != nil {
			return err
		}

		user.LatestLoggedInAt = ctx.MustGet("UtcNow").(time.Time)
		if err := u.ur.Save(ctx, *user); err != nil {
			return err
		}

		userAuth.Token = token
		userAuth.ExpiredAt = ctx.MustGet("UtcNow").(time.Time).AddDate(0, 0, 1)

		if err := u.uar.Save(ctx, *userAuth); err != nil {
			return err
		}

		return nil
	})

	return user, userAuth, err
}

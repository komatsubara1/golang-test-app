package use_case

import (
	"app/context"
	userEntity "app/domain/entity/user"
	user2 "app/domain/value/user"
	"app/infrastructure/repository/user"
	"app/lib"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type DebugUseCase struct {
}

// NewDebugUseCase デバッグユースケース生成
func NewDebugUseCase() *DebugUseCase {
	return &DebugUseCase{}
}

// SetTimeShift 時間遡行設定
func (u *DebugUseCase) SetTimeShift(ctx *gin.Context, timeShift int64) error {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)
	return lib.SetTimeShift(gctx, timeShift)
}

// SendPresent プレゼント付与
func (u *DebugUseCase) SendPresent(
	ctx *gin.Context,
	title string,
	params string,
	contentType uint,
	contentId uint64,
	amount uint64,
	arriveDate *time.Time,
	limitDate *time.Time,
) (*userEntity.UserPresent, error) {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)

	userPresentEntity := &userEntity.UserPresent{
		ID:          user2.NewPresentId(uuid.New()),
		UserId:      *gctx.UserId,
		Title:       title,
		Params:      params,
		ContentType: contentType,
		ContentId:   contentId,
		Amount:      amount,
		ArriveDate:  arriveDate,
		LimitDate:   limitDate,
	}

	userPresentRepository := user.NewUserPresentRepository()
	err := userPresentRepository.Save(ctx, *userPresentEntity)
	if err != nil {
		return nil, err
	}

	return userPresentEntity, nil
}

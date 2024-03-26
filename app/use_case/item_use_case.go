package use_case

import (
	"app/context"
	userEntity "app/domain/entity/user"
	"app/domain/enum/present"
	masterRepository "app/domain/repository/master"
	userRepository "app/domain/repository/user"
	"app/domain/service"
	masterValue "app/domain/value/master"
	userValue "app/domain/value/user"
	"app/infrastructure/repository/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type ItemUseCase struct {
	ur  userRepository.UserRepository
	uir userRepository.UserItemRepository
	upr userRepository.UserPresentRepository
	imr masterRepository.ItemMasterRepository
}

// NewItemUseCase アイテムユースケース生成
func NewItemUseCase(
	ur userRepository.UserRepository,
	uir userRepository.UserItemRepository,
	upr userRepository.UserPresentRepository,
	imr masterRepository.ItemMasterRepository) *ItemUseCase {
	return &ItemUseCase{ur, uir, upr, imr}
}

// Get アイテム取得
func (u *ItemUseCase) Get(ctx *gin.Context, itemId masterValue.ItemId) (*userEntity.UserItem, error) {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)
	userId := gctx.UserId

	item, err := u.uir.FindByUserIdAndItemId(ctx, *userId, itemId)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetAll アイテム取得
func (u *ItemUseCase) GetAll(ctx *gin.Context, exclusionZeroQuantity bool) (*userEntity.UserItems, error) {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)
	userId := gctx.UserId

	items, err := u.uir.FindByUserId(ctx, *userId)
	if err != nil {
		return nil, err
	}

	// 0以下を除外
	if exclusionZeroQuantity {
		items = service.UserItemService{}.FilterIsInStock(items)
	}

	upr := user.NewUserPresentRepository()
	_, _ = upr.FindByID(ctx, userValue.NewPresentId(uuid.New()))
	_, _ = upr.FindByUserIdAndArriveDateAndLimitDateAndReceivedAt(
		ctx,
		userValue.NewUserId(uuid.New()),
		ctx.MustGet("UtcNow").(time.Time),
		ctx.MustGet("UtcNow").(time.Time),
		ctx.MustGet("UtcNow").(time.Time),
	)

	return items, nil
}

// Gain アイテム付与
func (u *ItemUseCase) Gain(ctx *gin.Context, itemId masterValue.ItemId, count uint64) (*userEntity.UserItem, error) {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)
	userId := gctx.UserId

	var userItem *userEntity.UserItem = nil
	udctx := gctx.Udctx
	err := udctx.TransactionScope(func() error {

		// ユーザー情報取得
		user, err := u.ur.FindByID(ctx, *userId)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user not found. userId=%s", userId)
		}

		// アイテム情報取得
		item, err := u.imr.FindByID(ctx, itemId)
		if err != nil {
			return err
		}
		if item == nil {
			return fmt.Errorf("item not found. itemId=%d", itemId)
		}

		// ユーザーアイテム情報取得
		userItem, err = u.uir.FindByUserIdAndItemId(ctx, *userId, itemId)
		if err != nil {
			return err
		}
		if userItem != nil {
			userItem.Quantity = userItem.Quantity + count
		} else {
			userItem = service.UserItemService{}.Build(*userId, itemId, count)
		}

		var up *userEntity.UserPresent = nil
		var exceeded = service.UserItemService{}.IsExceedMaxCount(userItem, item)
		if exceeded {
			up = service.UserPresentService{}.Build(
				*userId,
				// TODO:
				"test",
				"",
				uint(present.ContentTypeItem),
				itemId.Value(),
				count,
				nil,
			)
		}

		if err := u.uir.Save(ctx, *userItem); err != nil {
			return err
		}

		if up != nil {
			if err := u.upr.Save(ctx, *up); err != nil {
				return err
			}
		}

		return nil
	})

	return userItem, err
}

// Use アイテム使用
func (u *ItemUseCase) Use(ctx *gin.Context, itemId masterValue.ItemId, count uint64) (*userEntity.User, *userEntity.UserItem, error) {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)
	userId := gctx.UserId

	var user *userEntity.User = nil
	var userItem *userEntity.UserItem = nil
	udctx := gctx.Udctx
	err := udctx.TransactionScope(func() error {
		// ユーザー情報取得
		user, err := u.ur.FindByID(ctx, *userId)
		if err != nil {
			return err
		}

		// ユーザーアイテム情報取得
		userItem, err := u.uir.FindByUserIdAndItemId(ctx, *userId, itemId)
		if err != nil {
			return err
		}

		// 所持数検証
		if userItem.Quantity < count {
			return fmt.Errorf("item quantity not enough for use. item_id=%d, use_count=%d, quantity=%d", itemId, count, userItem.Quantity)
		}

		item, err := u.imr.FindByID(ctx, itemId)
		if err != nil {
			return err
		}

		// ユーザーアイテム所持数減算
		userItem.Quantity = userItem.Quantity - count

		if err := u.uir.Save(ctx, *userItem); err != nil {
			return err
		}

		// 効果発動
		// スタミナ回復のみ考慮
		switch item.EffectType {
		case 1: // TODO: enum_item.ItemEffectTypeRecoverStamina:
			user = service.UserItemService{}.RecoverStamina(user, item, ctx.MustGet("UtcNow").(time.Time))
			if err := u.ur.Save(ctx, *user); err != nil {
				return err
			}
		default:
			return fmt.Errorf("undefined EffectType. EffectType=%d", item.EffectType)
		}

		return nil
	})

	return user, userItem, err
}

// Sell アイテム売却
func (u *ItemUseCase) Sell(ctx *gin.Context, itemId masterValue.ItemId, count uint64) (*userEntity.User, *userEntity.UserItem, error) {
	gctx := ctx.MustGet("GameContext").(*context.GameContext)
	userId := gctx.UserId

	var user *userEntity.User = nil
	var userItem *userEntity.UserItem = nil
	udctx := gctx.Udctx
	err := udctx.TransactionScope(func() error {
		// ユーザーアイテム情報取得
		userItem, err := u.uir.FindByUserIdAndItemId(ctx, *userId, itemId)
		if err != nil {
			return err
		}

		// 所持数検証
		if userItem.Quantity < count {
			return fmt.Errorf("item quantity not enough for sell. sellCount=%d, quantity=%d", count, userItem.Quantity)
		}

		// ユーザーアイテム所持数減算
		userItem.Quantity = userItem.Quantity - count

		item, err := u.imr.FindByID(ctx, itemId)
		if err != nil {
			return err
		}

		user, err = u.ur.FindByID(ctx, *userId)
		if err != nil {
			return err
		}

		if err = u.uir.Save(ctx, *userItem); err != nil {
			return err
		}

		// ソフトカレンシー付与
		if item.SellCoin != 0 {
			addCount := item.SellCoin * count
			user.Coin = user.Coin + addCount
			if err := u.ur.Save(ctx, *user); err != nil {
				return err
			}
		}

		return nil
	})

	return user, userItem, err
}

package user

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/exp/slog"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"app/context"
	userentity "app/domain/entity/user"
	userrepository "app/domain/repository/user"
	mastervalue "app/domain/value/master"
	uservalue "app/domain/value/user"
)

type userItemRepository struct {
}

func NewUserItemRepository() userrepository.UserItemRepository {
	return &userItemRepository{}
}

func (r *userItemRepository) FindByUserIdAndItemId(
	ctx *gin.Context, userId uservalue.UserId, itemId mastervalue.ItemId,
) (*userentity.UserItem, error) {
	entity := &userentity.UserItem{}

	log.Printf("%s: %d", userId.Value(), itemId.Value())

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("user_id, item_id, quantity").
		From(entity.TableName()).
		Where(squirrel.Eq{"user_id": userId.Value(), "item_id": itemId.Value()})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info("execute query", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	err = builder.RunWith(udctx.Dc.Db).
		QueryRowContext(ctx).
		Scan(&entity.UserId, &entity.ItemId, &entity.Quantity)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("user_items no records.", "userId", userId, "itemId", itemId)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// FindByUserId ユーザー所持アイテム情報取得
func (r *userItemRepository) FindByUserId(ctx *gin.Context, userId uservalue.UserId) (*userentity.UserItems, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("user_id, item_id, quantity").From("user_item").Where(squirrel.Eq{"user_id": userId.Value().String()})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info("execute query", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	rows, err := builder.RunWith(udctx.Dc.Db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	userItems := userentity.UserItems{}
	for rows.Next() {
		userItem := &userentity.UserItem{}
		if err := rows.Scan(&userItem.UserId, &userItem.ItemId, &userItem.Quantity); err != nil {
			return nil, err
		}

		userItems = append(userItems, *userItem)
	}

	return &userItems, nil
}

// Save 更新
func (r *userItemRepository) Save(ctx *gin.Context, entity userentity.UserItem) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert(entity.TableName()).Columns("user_id, item_id, quantity").Values(entity.UserId.Value(), entity.ItemId.Value(), entity.Quantity).
		Suffix("ON DUPLICATE KEY UPDATE quantity = ?", entity.Quantity)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	slog.Info("execute query", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	res, err := builder.RunWith(udctx.Dc.Db).ExecContext(ctx)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		slog.Info("no affected execute query.", query, args, err)
	}

	return nil
}

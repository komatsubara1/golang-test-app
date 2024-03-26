package user

import (
	"database/sql"
	"errors"

	"golang.org/x/exp/slog"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"

	"app/context"
	userentity "app/domain/entity/user"
	userrepository "app/domain/repository/user"
	uservalue "app/domain/value/user"
)

type userAuthRepository struct {
}

func NewUserAuthRepository() userrepository.UserAuthRepository {
	return &userAuthRepository{}
}

// FindByUserId 認証情報取得
func (r *userAuthRepository) FindByUserId(ctx *gin.Context, userId uservalue.UserId) (*userentity.UserAuth, error) {
	entity := &userentity.UserAuth{}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("user_id, token").
		From(entity.TableName()).
		Where(squirrel.Eq{"user_id": userId.Value()})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info("execute query.", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	err = builder.RunWith(udctx.Dc.Db).
		QueryRowContext(ctx).
		Scan(&entity.UserId, &entity.Token)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("user_auths no records.", "userId", userId)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// Save 認証情報登録
func (r *userAuthRepository) Save(ctx *gin.Context, entity userentity.UserAuth) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert(entity.TableName()).
		Columns("user_id, token, expired_at").
		Values(entity.UserId.Value(), entity.Token, entity.ExpiredAt).
		Suffix("ON DUPLICATE KEY UPDATE token = ?, expired_at = ?", entity.Token, entity.ExpiredAt)

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
		slog.Info("RowsAffected zero.")
	}

	return nil
}

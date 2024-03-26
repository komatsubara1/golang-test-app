package user

import (
	"database/sql"
	"errors"
	"golang.org/x/exp/slog"

	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"app/context"
	userentity "app/domain/entity/user"
	userrepository "app/domain/repository/user"
	uservalue "app/domain/value/user"
)

type userRepository struct {
}

func NewUserRepository() userrepository.UserRepository {
	return &userRepository{}
}

// FindByID ユーザー取得
func (r *userRepository) FindByID(ctx *gin.Context, id uservalue.UserId) (*userentity.User, error) {
	user := &userentity.User{}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("id, name, stamina, stamina_latest_updated_at, coin, latest_logged_in_at").
		From(user.TableName()).Where(squirrel.Eq{"id": id.Value().String()})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info("execute query", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	err = builder.RunWith(udctx.Dc.Db).
		QueryRowContext(ctx).
		Scan(&user.ID, &user.Name, &user.Stamina, &user.StaminaLatestUpdatedAt,
			&user.Coin, &user.LatestLoggedInAt)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("users no records.", "id", id)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Save(ctx *gin.Context, user userentity.User) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert(user.TableName()).Columns("id, name, stamina, stamina_latest_updated_at, coin, latest_logged_in_at").
		Values(user.ID.Value(),
			user.Name,
			user.Stamina,
			user.StaminaLatestUpdatedAt,
			user.Coin,
			user.LatestLoggedInAt).
		Suffix("ON DUPLICATE KEY UPDATE stamina = ?, stamina_latest_updated_at = ?, coin = ?, latest_logged_in_at = ?",
			user.Stamina, user.StaminaLatestUpdatedAt, user.Coin, user.LatestLoggedInAt)

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
		slog.Info("no affected execute query.", "query", query, "args", args, "error", err)
	}

	return nil
}

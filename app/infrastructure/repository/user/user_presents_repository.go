package user

import (
	"app/context"
	userentity "app/domain/entity/user"
	userrepository "app/domain/repository/user"
	uservalue "app/domain/value/user"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"time"
)

type userPresentRepository struct {
}

func NewUserPresentRepository() userrepository.UserPresentRepository {
	return &userPresentRepository{}
}

func (r *userPresentRepository) FindByID(ctx *gin.Context, id uservalue.PresentId) (*userentity.UserPresent, error) {
	entity := &userentity.UserPresent{}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("id, user_id, title, params, content_type, content_id, amount, arrive_date, limit_date, received_at").
		From(entity.TableName()).
		Where(squirrel.Eq{"id": id.Value()}).
		Limit(100)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info("execute query.", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	err = builder.RunWith(udctx.Dc.Db).
		QueryRowContext(ctx).
		Scan(&entity.ID, &entity.UserId, &entity.Title, &entity.Params, &entity.ContentType, &entity.ContentId, &entity.Amount, &entity.ArriveDate, &entity.LimitDate, &entity.ReceivedAt)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("user_presents no records.", "id", id)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// FindByUserIdAndArriveDateAndLimitDateAndReceivedAt TODO: FindAvailable
func (r *userPresentRepository) FindByUserIdAndArriveDateAndLimitDateAndReceivedAt(
	ctx *gin.Context, userId uservalue.UserId, arriveDate time.Time, limitDate time.Time, receivedAt time.Time,
) (*userentity.UserPresents, error) {
	entity := &userentity.UserPresent{}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select("id, user_id, title, params, content_type, content_id, amount, arrive_date, limit_date, received_at").
		From(entity.TableName()).
		Where(squirrel.Eq{"user_id": userId.Value()}).
		Where(squirrel.GtOrEq{"arrive_date": ctx.MustGet("UtcNow").(time.Time)}).
		Where(squirrel.Or{squirrel.Eq{"limit_date": nil}, squirrel.GtOrEq{"limit_date": ctx.MustGet("UtcNow").(time.Time).UTC()}}).
		Where(squirrel.Eq{"received_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	slog.Info("execute query.", "query", query, "args", args)

	udctx := ctx.MustGet("GameContext").(*context.GameContext).Udctx
	udctx.Connect()
	rows, err := builder.RunWith(udctx.Dc.Db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	userPresents := userentity.UserPresents{}
	for rows.Next() {
		userPresent := &userentity.UserPresent{}
		err := rows.Scan(&userPresent.ID, &userPresent.UserId, &userPresent.Title, &userPresent.Params, &userPresent.ContentType, &userPresent.ContentId, &userPresent.Amount, &userPresent.ArriveDate, &userPresent.LimitDate, &userPresent.ReceivedAt)
		if err != nil {
			return nil, err
		}

		userPresents = append(userPresents, *userPresent)
	}

	return &userPresents, nil
}

func (r *userPresentRepository) Save(ctx *gin.Context, entity userentity.UserPresent) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert(entity.TableName()).
		Columns("id, user_id, title, params, content_type, content_id, amount, arrive_date, limit_date, received_at").
		Values(entity.ID.Value(), entity.UserId.Value(), entity.Title, entity.Params, entity.ContentType, entity.ContentId, entity.Amount, entity.ArriveDate, entity.LimitDate, entity.ReceivedAt).
		Suffix("")

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

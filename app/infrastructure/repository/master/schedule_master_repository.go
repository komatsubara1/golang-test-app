package master

import (
	"app/context"
	"fmt"
	"gorm.io/gorm/clause"

	masterentity "app/domain/entity/master"
	masterrepository "app/domain/repository/master"
	mastervalue "app/domain/value/master"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type scheduleMasterRepository struct {
}

func NewScheduleMasterRepository() masterrepository.ScheduleMasterRepository {
	return &scheduleMasterRepository{}
}

func (r *scheduleMasterRepository) FindByID(ctx *gin.Context, id mastervalue.ScheduleId) (*masterentity.ScheduleMaster, error) {
	// TODO: user_redis get

	i, err := r.findByIdForDB(ctx, id)
	if err != nil {
		return nil, err
	}
	if i == nil {
		slog.Info("schedule master not found.", "itemId", id)
		return nil, fmt.Errorf("schedule master not found. id=%d", id)
	}

	// TODO: user_redis set

	return i, nil
}

func (r *scheduleMasterRepository) findByIdForDB(ctx *gin.Context, id mastervalue.ScheduleId) (*masterentity.ScheduleMaster, error) {
	schedule := &masterentity.ScheduleMaster{}
	mdctx := ctx.MustGet("GameContext").(*context.GameContext).Mdctx
	mdctx.Connect()
	err := mdctx.Dc.Db.Where("id = ?", id.Value()).First(&schedule).Error
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *scheduleMasterRepository) Save(ctx *gin.Context, entity masterentity.ScheduleMaster) error {
	mdctx := ctx.MustGet("GameContext").(*context.GameContext).Mdctx
	mdctx.Connect()
	err := mdctx.Dc.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"start_at",
			"end_at",
			"close_at",
		}),
	}).Create(entity).Error
	if err != nil {
		return err
	}

	// TODO: Redis Set

	return nil
}

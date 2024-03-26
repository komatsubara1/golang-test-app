package master

import (
	"app/context"
	masterentity "app/domain/entity/master"
	masterrepository "app/domain/repository/master"
	mastervalue "app/domain/value/master"
	"app/lib"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

const ItemMasterRedisKey = "item_masters"

type itemMasterRepository struct {
}

func NewItemMasterRepository() masterrepository.ItemMasterRepository {
	return &itemMasterRepository{}
}

func (r *itemMasterRepository) FindByID(ctx *gin.Context, id mastervalue.ItemId) (*masterentity.ItemMaster, error) {
	mcctx := ctx.MustGet("GameContext").(*context.GameContext).Mcctxt
	cache := mcctx.Connect()
	ic := lib.HGet[masterentity.ItemMaster](cache, ItemMasterRedisKey, id)
	if ic != nil {
		return ic, nil
	}

	i, err := r.findByIdForDB(ctx, id)
	if err != nil {
		return nil, err
	}
	if i == nil {
		slog.Info("item master not found.", "itemId", id)
		return nil, nil
	}

	err = lib.HSet[masterentity.ItemMaster](cache, ItemMasterRedisKey, id, i)
	if err != nil {
		slog.Warn("item master not cache saved.", "itemId", id)
		return nil, nil
	}

	return i, nil
}

func (r *itemMasterRepository) findByIdForDB(ctx *gin.Context, id mastervalue.ItemId) (*masterentity.ItemMaster, error) {
	entity := &masterentity.ItemMaster{}
	mdctx := ctx.MustGet("GameContext").(*context.GameContext).Mdctx
	mdctx.Connect()
	err := mdctx.Dc.Db.Where("id = ?", id.Value()).First(entity).Error
	return entity, err
}

func (r *itemMasterRepository) Save(ctx *gin.Context, entity masterentity.ItemMaster) error {
	mdctx := ctx.MustGet("GameContext").(*context.GameContext).Mdctx
	mdctx.Connect()
	err := mdctx.Dc.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"type",
			"sell_coin",
			"effect_type",
			"effect_value",
			"schedule_id",
			"max_count",
			"max_view_count",
		}),
	}).Create(entity).Error
	if err != nil {
		return err
	}

	mcctx := ctx.MustGet("GameContext").(*context.GameContext).Mcctxt
	cache := mcctx.Connect()
	err = lib.HSet[masterentity.ItemMaster](cache, ItemMasterRedisKey, entity.ID, &entity)
	if err != nil {
		slog.Warn("item master not cache saved.", "itemId", entity.ID)
	}

	return nil
}

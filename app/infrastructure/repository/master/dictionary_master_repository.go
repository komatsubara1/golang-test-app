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

const DictionaryMasterRedisKey = "dictionary_masters"

type dictionaryMasterRepository struct {
}

func NewDictionaryMasterRepository() masterrepository.DictionaryMasterRepository {
	return &dictionaryMasterRepository{}
}

func (r *dictionaryMasterRepository) FindByKey(
	ctx *gin.Context,
	key mastervalue.DictionaryKey,
) (*masterentity.DictionaryMaster, error) {
	mcctx := ctx.MustGet("GameContext").(*context.GameContext).Mcctxt
	cache := mcctx.Connect()
	ic := lib.HGet[masterentity.DictionaryMaster](cache, DictionaryMasterRedisKey, key)
	if ic != nil {
		return ic, nil
	}

	i, err := r.findByIdForDB(ctx, key)
	if err != nil {
		return nil, err
	}
	if i == nil {
		slog.Info("dictionary master not found.", "DictionaryKey", key)
		return nil, nil
	}

	err = lib.HSet[masterentity.DictionaryMaster](cache, DictionaryMasterRedisKey, key, i)
	if err != nil {
		slog.Warn("dictionary master not cache saved.", "DictionaryKey", key)
		return nil, nil
	}

	return i, nil
}

func (r *dictionaryMasterRepository) findByIdForDB(
	ctx *gin.Context,
	key mastervalue.DictionaryKey,
) (*masterentity.DictionaryMaster, error) {
	entity := &masterentity.DictionaryMaster{}
	mdctx := ctx.MustGet("GameContext").(*context.GameContext).Mdctx
	mdctx.Connect()
	err := mdctx.Dc.Db.Where("id = ?", key.Value()).First(entity).Error
	return entity, err
}

func (r *dictionaryMasterRepository) Save(
	ctx *gin.Context,
	entity masterentity.DictionaryMaster,
) error {
	mdctx := ctx.MustGet("GameContext").(*context.GameContext).Mdctx
	mdctx.Connect()
	err := mdctx.Dc.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "Key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"en",
			"ja",
		}),
	}).Create(entity).Error
	if err != nil {
		return err
	}

	mcctx := ctx.MustGet("GameContext").(*context.GameContext).Mcctxt
	cache := mcctx.Connect()
	err = lib.HSet[masterentity.DictionaryMaster](cache, DictionaryMasterRedisKey, entity.Key, &entity)
	if err != nil {
		slog.Warn("dictionary master not cache saved.", "DictionaryKey", entity.Key)
	}

	return nil
}

package context

import (
	"app/domain/value/user"
	"time"
)

type GameContext struct {
	UserId *user.UserId

	Udctx  *UserDbContext
	Ucctx  *UserCacheContext
	Mdctx  *MasterDbContext
	Mcctxt *MasterCacheContext

	UtcNow time.Time
}

func NewGameContext(
	udctx *UserDbContext,
	ucctx *UserCacheContext,
	mdctx *MasterDbContext,
	mcctxt *MasterCacheContext,
	utcNow time.Time,
) *GameContext {
	return &GameContext{UserId: nil, Udctx: udctx, Ucctx: ucctx, Mdctx: mdctx, Mcctxt: mcctxt, UtcNow: utcNow}
}

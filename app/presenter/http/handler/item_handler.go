package handler

import (
	"app/domain/enum/error"
	request "app/presenter/http/handler/request/item"
	response "app/presenter/http/handler/response/item"
	"app/use_case"
	"net/http"

	mastervalue "app/domain/value/master"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type ItemHandler struct {
	u *use_case.ItemUseCase
}

func NewItemHandler(u *use_case.ItemUseCase) *ItemHandler {
	return &ItemHandler{u}
}

// Get
// @Summary 対象のユーザー所持アイテム情報を取得
// @Tags Item
// @Accept  json
// @Produce  json
// @Param title body request.ItemGetRequest true "アイテム取得リクエスト"
// @Success 200 {object} response.ItemGetResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/get [POST]
func (h *ItemHandler) Get(ctx *gin.Context) {
	// validate req
	var req request.ItemGetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("validate error.", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemId := req.ItemId

	userItem, err := h.u.Get(ctx, itemId)
	if err != nil {
		slog.Error("Item Get error.", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewItemGetResponse("ok", error.ErrorCodeNone, userItem)
	ctx.JSON(http.StatusOK, r)
}

// GetAll
// @Summary ユーザー所持アイテム情報を取得
// @Tags Item
// @Accept  json
// @Produce  json
// @Param title body request.ItemGetAllRequest true "アイテム取得リクエスト"
// @Success 200 {object} response.ItemGetAllResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/get_all [POST]
func (h *ItemHandler) GetAll(ctx *gin.Context) {
	// validate req
	var req request.ItemGetAllRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userItems, err := h.u.GetAll(ctx, req.ExclusionZeroQuantity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewItemGetAllResponse("ok", error.ErrorCodeNone, userItems)
	ctx.JSON(http.StatusOK, r.ToJson())
}

// Gain
// @Summary アイテムをユーザーに付与
// @Tags Item
// @Accept  json
// @Produce  json
// @Param title body request.ItemGainRequest true "アイテム付与リクエスト"
// @Success 200 {object} response.ItemGainResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/gain [POST]
func (h *ItemHandler) Gain(ctx *gin.Context) {
	// validate req
	var req request.ItemGainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemId := mastervalue.NewItemId(req.ItemId)
	count := req.Count

	userItem, err := h.u.Gain(ctx, itemId, count)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewItemGainResponse("ok", error.ErrorCodeNone, userItem)
	ctx.JSON(http.StatusOK, r.ToJson())
}

// Use
// @Summary アイテムを使用
// @Tags Item
// @Accept  json
// @Produce  json
// @Param title body request.ItemUseRequest true "アイテム使用リクエスト"
// @Success 200 {object} response.ItemUseResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/gain [POST]
func (h *ItemHandler) Use(ctx *gin.Context) {
	// validate req
	var req request.ItemUseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemId := req.ItemId
	count := req.Count

	user, userItem, err := h.u.Use(ctx, itemId, count)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewItemUseResponse("ok", error.ErrorCodeNone, user, userItem)
	ctx.JSON(http.StatusOK, r.ToJson())
}

// Sell
// @Summary アイテムを売却
// @Tags Item
// @Accept  json
// @Produce  json
// @Param title body request.ItemSellRequest true "アイテム売却リクエスト"
// @Success 200 {object} response.ItemSellResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/gain [POST]
func (h *ItemHandler) Sell(ctx *gin.Context) {
	// validate req
	var req request.ItemSellRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemId := req.ItemId
	count := req.Count

	user, userItem, err := h.u.Sell(ctx, itemId, count)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewItemSellResponse("ok", error.ErrorCodeNone, user, userItem)
	ctx.JSON(http.StatusOK, r.ToJson())
}

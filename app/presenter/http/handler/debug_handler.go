package handler

import (
	"app/domain/enum/error"
	request "app/presenter/http/handler/request/debug"
	response "app/presenter/http/handler/response/debug"
	"app/use_case"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type DebugHandler struct {
	u *use_case.DebugUseCase
}

func NewDebugHandler(u *use_case.DebugUseCase) *DebugHandler {
	return &DebugHandler{u}
}

// SetTimeShift
// @Summary 時間遡行の差分を設定
// @Tags Debug
// @Accept  json
// @Produce  json
// @Param title body request.DebugSetTimeShiftRequest true "時間遡行設定リクエスト"
// @Success 200 {object} response.DebugSetTimeShiftResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/get [POST]
func (h *DebugHandler) SetTimeShift(ctx *gin.Context) {
	// validate req
	var req request.DebugSetTimeShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("validate error.", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeShift := req.TimeShift

	err := h.u.SetTimeShift(ctx, timeShift)
	if err != nil {
		slog.Error("Debug SetTimeShift error.", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewDebugSetTimeShiftResponse("ok", error.ErrorCodeNone)
	ctx.JSON(http.StatusOK, r)
}

// SendPresent
// @Summary プレゼント付与
// @Tags Debug
// @Accept  json
// @Produce  json
// @Param title body request.DebugSendPresentRequest true "プレゼント付与リクエスト"
// @Success 200 {object} response.DebugSendPresentResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /item/get [POST]
func (h *DebugHandler) SendPresent(ctx *gin.Context) {
	// validate req
	var req request.DebugSendPresentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("validate error.", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title := req.Title
	params := req.Params
	contentType := req.ContentType
	contentId := req.ContentId
	amount := req.Amount
	arriveDate := req.ArriveDate
	limitDate := req.LimitDate

	userPresent, err := h.u.SendPresent(ctx, title, params, contentType, contentId, amount, arriveDate, limitDate)
	if err != nil {
		slog.Error("Debug SetTimeShift error.", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewDebugSendPresentResponse("ok", error.ErrorCodeNone, *userPresent)
	ctx.JSON(http.StatusOK, r)
}

package handler

import (
	"app/context"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/domain/enum/error"
	request "app/presenter/http/handler/request/user"
	response "app/presenter/http/handler/response/user"
	"app/use_case"
)

type UserHandler struct {
	u *use_case.UserUseCase
}

func NewUserHandler(u *use_case.UserUseCase) *UserHandler {
	return &UserHandler{u}
}

// Get
// @Summary ユーザーを取得
// @Tags User
// @Accept  json
// @Produce  json
// @Param title body request.UserGetRequest true "ユーザー取得リクエスト"
// @Success 200 {object} response.UserGetResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /user/get [POST]
func (h UserHandler) Get(ctx *gin.Context) {
	userId := ctx.MustGet("GameContext").(*context.GameContext).UserId

	user, err := h.u.Get(ctx, *userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	r := response.NewUserGetResponse("ok", error.ErrorCodeNone, user)
	ctx.JSON(http.StatusOK, r.ToJson())
}

// Create
// @Summary ユーザーを作成
// @Tags User
// @Accept  json
// @Produce  json
// @Param title body request.UserCreateRequest true "ユーザー作成リクエスト"
// @Success 200 {object} response.UserCreateResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /user/create [POST]
func (h UserHandler) Create(ctx *gin.Context) {
	// validate req
	var req request.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, userAuthEntity, err := h.u.Create(ctx, req.UserName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", userAuthEntity.Token)

	// response
	r := response.NewUserCreateResponse("ok", error.ErrorCodeNone, user)
	ctx.JSON(http.StatusOK, r.ToJson())
}

// Login
// @Summary ユーザーログイン
// @Tags User
// @Accept  json
// @Produce  json
// @Param title body request.UserLoginRequest true "ユーザーログインリクエスト"
// @Success 200 {object} response.UserLoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /user/login [POST]
func (h UserHandler) Login(ctx *gin.Context) {
	// validate req
	var req request.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, userAuth, err := h.u.Login(ctx, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", userAuth.Token)

	// response
	r := response.NewUserLoginResponse("ok", error.ErrorCodeNone, user)
	ctx.JSON(http.StatusOK, r.ToJson())
}

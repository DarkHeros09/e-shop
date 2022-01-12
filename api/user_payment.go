package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserPaymentRequest struct {
	UserID      int64     `json:"user_id" binding:"required"`
	PaymentType string    `json:"payment_type" binding:"required"`
	Provider    string    `json:"provider" binding:"required"`
	AccountNo   int32     `json:"account_no" binding:"required"`
	Expiry      time.Time `json:"expiry" binding:"required"`
}

func (server *Server) createUserPayment(ctx *gin.Context) {
	var req createUserPaymentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateUserPaymentParams{
		UserID:      req.UserID,
		PaymentType: req.PaymentType,
		Provider:    req.Provider,
		AccountNo:   req.AccountNo,
		Expiry:      req.Expiry,
	}

	userPayment, err := server.store.CreateUserPayment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userPayment)
}

type getUserPaymentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserPayment(ctx *gin.Context) {
	var req getUserPaymentRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	userPayment, err := server.store.GetUserPayment(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userPayment)
}

type listUserPaymentsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUserPayments(ctx *gin.Context) {
	var req listUserPaymentsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.ListUserPaymentsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	userPayments, err := server.store.ListUserPayments(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userPayments)
}

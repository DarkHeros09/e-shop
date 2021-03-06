package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserPaymentRequest struct {
	UserID      int64     `json:"user_id" binding:"required,min=1"`
	PaymentType string    `json:"payment_type" binding:"required"`
	Provider    string    `json:"provider" binding:"required"`
	AccountNo   int32     `json:"account_no" binding:"required,numeric"`
	Expiry      time.Time `json:"expiry" binding:"required"`
}

func (server *Server) createUserPayment(ctx *gin.Context) {
	var req createUserPaymentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.CreateUserPaymentParams{
		UserID:      authPayload.UserID,
		PaymentType: req.PaymentType,
		Provider:    req.Provider,
		AccountNo:   req.AccountNo,
		Expiry:      req.Expiry,
	}

	userPayment, err := server.store.CreateUserPayment(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.GetUserPaymentParams{
		ID:     req.ID,
		UserID: authPayload.UserID,
	}

	userPayment, err := server.store.GetUserPayment(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.ListUserPaymentsParams{
		UserID: authPayload.UserID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	userPayments, err := server.store.ListUserPayments(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userPayments)
}

type updateUserPaymentRequest struct {
	ID          int64  `json:"id" binding:"required,min=1"`
	PaymentType string `json:"payment_type" binding:"required"`
}

func (server *Server) updateUserPayment(ctx *gin.Context) {
	var req updateUserPaymentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.UpdateUserPaymentParams{
		ID:          req.ID,
		UserID:      authPayload.UserID,
		PaymentType: req.PaymentType,
	}

	userPayment, err := server.store.UpdateUserPayment(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userPayment)
}

type deleteUserPaymentRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (server *Server) deleteUserPayment(ctx *gin.Context) {
	var req deleteUserPaymentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.DeleteUserPaymentParams{
		ID:     req.ID,
		UserID: authPayload.UserID,
	}

	err := server.store.DeleteUserPayment(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		} else if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type getPaymentDetailRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPaymentDetail(ctx *gin.Context) {
	var req getPaymentDetailRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.GetPaymentDetailParams{
		ID:     req.ID,
		UserID: authPayload.UserID,
	}

	paymentDetail, err := server.store.GetPaymentDetail(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, paymentDetail)
}

type listPaymentDetailsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listPaymentDetails(ctx *gin.Context) {
	var req listPaymentDetailsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.ListPaymentDetailsParams{
		UserID: authPayload.UserID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	paymentDetails, err := server.store.ListPaymentDetails(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, paymentDetails)
}

type updatePaymentDetailRequest struct {
	ID       int64  `json:"id" binding:"required,min=1"`
	OrderID  int64  `json:"order_id"`
	Amount   int32  `json:"amount"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

func (server *Server) updatePaymentDetail(ctx *gin.Context) {
	var req updatePaymentDetailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.UpdatePaymentDetailParams{
		ID:       req.ID,
		UserID:   authPayload.UserID,
		OrderID:  req.OrderID,
		Amount:   req.Amount,
		Provider: req.Provider,
		Status:   req.Status,
	}

	paymentDetail, err := server.store.UpdatePaymentDetail(ctx, arg)
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

	ctx.JSON(http.StatusOK, paymentDetail)
}

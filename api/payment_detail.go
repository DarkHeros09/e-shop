package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
)

// type createPaymentDetailRequest struct {
// 	OrderID  int64  `json:"order_id" binding:"required,min=1"`
// 	Amount   int32  `json:"amount" binding:"required"`
// 	Provider string `json:"provider" binding:"required"`
// 	Status   string `json:"status" binding:"required"`
// }

// func (server *Server) createPaymentDetail(ctx *gin.Context) {
// 	var req createPaymentDetailRequest

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	orderDetail, err := server.store.GetOrderDetail(ctx, req.OrderID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
// 	if orderDetail.UserID != authPayload.UserID {
// 		err := errors.New("account deosn't belong to the authenticated user")
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
// 		return
// 	}

// 	arg := db.CreatePaymentDetailParams{
// 		OrderID:  req.OrderID,
// 		Amount:   req.Amount,
// 		Provider: req.Provider,
// 		Status:   req.Status,
// 	}

// 	paymentDetail, err := server.store.CreatePaymentDetail(ctx, arg)
// 	if err != nil {
// 		if pqErr, ok := err.(*pq.Error); ok {
// 			switch pqErr.Code.Name() {
// 			case "unique_violation":
// 				ctx.JSON(http.StatusForbidden, errorResponse(err))
// 				return
// 			}
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, paymentDetail)
// }

type getPaymentDetailRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getPaymentDetail(ctx *gin.Context) {
	var req getPaymentDetailRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	paymentDetail, err := server.store.GetPaymentDetail(ctx, req.ID)
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

	arg := db.ListPaymentDetailsParams{
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

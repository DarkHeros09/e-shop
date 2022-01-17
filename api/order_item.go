package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createOrderItemRequest struct {
	OrderID   int64 `json:"order_id" binding:"required"`
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int32 `json:"quantity" binding:"required"`
}

func (server *Server) createOrderItem(ctx *gin.Context) {
	var req createOrderItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateOrderItemParams{
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	orderItem, err := server.store.CreateOrderItem(ctx, arg)
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

	ctx.JSON(http.StatusOK, orderItem)
}

type getOrderItemByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrderItemByID(ctx *gin.Context) {
	var req getOrderItemByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderItem, err := server.store.GetOrderItemByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, orderItem)
}

type getOrderItemByOrderDetailIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrderItemByOrderDetailID(ctx *gin.Context) {
	var req getOrderItemByOrderDetailIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderItem, err := server.store.GetOrderItemByOrderDetailID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, orderItem)
}

type listOrderItemsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listOrderItems(ctx *gin.Context) {
	var req listOrderItemsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListOrderItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	orderItems, err := server.store.ListOrderItems(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, orderItems)
}

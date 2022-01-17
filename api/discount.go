package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createDiscountRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	DiscountPercent string `json:"discount_percent" binding:"required"`
}

func (server *Server) createDiscount(ctx *gin.Context) {
	var req createDiscountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateDiscountParams{
		Name:            req.Name,
		Description:     req.Description,
		DiscountPercent: req.DiscountPercent,
	}

	discount, err := server.store.CreateDiscount(ctx, arg)
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

	ctx.JSON(http.StatusOK, discount)
}

type getDiscountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getDiscount(ctx *gin.Context) {
	var req getDiscountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	discount, err := server.store.GetDiscount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, discount)
}

type listDiscountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listDiscount(ctx *gin.Context) {
	var req listDiscountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListDiscountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	discounts, err := server.store.ListDiscounts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, discounts)
}

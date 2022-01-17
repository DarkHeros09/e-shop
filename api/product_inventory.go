package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createProductInventoryRequest struct {
	Quantity int32 `json:"quantity" binding:"required"`
}

func (server *Server) createInventory(ctx *gin.Context) {
	var req createProductInventoryRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := req.Quantity

	inventory, err := server.store.CreateProductInventory(ctx, arg)
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

	ctx.JSON(http.StatusOK, inventory)
}

type getProductInventoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getInventory(ctx *gin.Context) {
	var req getProductInventoryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	inventory, err := server.store.GetProductInventory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, inventory)
}

type listProductInventoriesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listInventories(ctx *gin.Context) {
	var req listProductInventoriesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductInventoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	inventories, err := server.store.ListProductInventories(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, inventories)
}

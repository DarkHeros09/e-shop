package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getDiscountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getDiscount(ctx *gin.Context) {
	var req getDiscountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	discount, err := server.store.GetDiscount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, discount)
}

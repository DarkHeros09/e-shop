package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createShoppingSessionRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Total  string `json:"total" binding:"required"`
}

func (server *Server) createShoppingSession(ctx *gin.Context) {
	var req createShoppingSessionRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateShoppingSessionParams{
		UserID: req.UserID,
		Total:  req.Total,
	}

	shoppingSession, err := server.store.CreateShoppingSession(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	ctx.JSON(http.StatusOK, shoppingSession)
}

type getShoppingSessionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getShoppingSession(ctx *gin.Context) {
	var req getShoppingSessionRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	shoppingSession, err := server.store.GetShoppingSession(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, shoppingSession)
}
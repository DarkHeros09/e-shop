package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createShoppingSessionRequest struct {
	UserID int64  `json:"user_id" binding:"required,min=1"`
	Total  string `json:"total" binding:"required,min=0"`
}

func (server *Server) createShoppingSession(ctx *gin.Context) {
	var req createShoppingSessionRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.CreateShoppingSessionParams{
		UserID: authPayload.UserID,
		Total:  req.Total,
	}

	shoppingSession, err := server.store.CreateShoppingSession(ctx, arg)
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

	ctx.JSON(http.StatusOK, shoppingSession)
}

type getShoppingSessionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getShoppingSession(ctx *gin.Context) {
	var req getShoppingSessionRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shoppingSession, err := server.store.GetShoppingSession(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	if shoppingSession.UserID != authPayload.UserID {
		err := errors.New("account deosn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, shoppingSession)
}

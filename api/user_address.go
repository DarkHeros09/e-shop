package api

import (
	"database/sql"
	"net/http"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserAddressRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	AddressLine string `json:"address_line" binding:"required"`
	City        string `json:"city" binding:"required"`
	Telephone   int32  `json:"telephone" binding:"required"`
}

func (server *Server) createUserAddress(ctx *gin.Context) {
	var req createUserAddressRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserAddressParams{
		UserID:      req.UserID,
		AddressLine: req.AddressLine,
		City:        req.City,
		Telephone:   req.Telephone,
	}

	userAddress, err := server.store.CreateUserAddress(ctx, arg)
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

	ctx.JSON(http.StatusOK, userAddress)
}

type getUserAddressRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserAddress(ctx *gin.Context) {
	var req getUserAddressRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userAddress, err := server.store.GetUserAddress(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userAddress)
}

type listUserAddressesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUserAddresses(ctx *gin.Context) {
	var req listUserAddressesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserAddressesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	userAddresses, err := server.store.ListUserAddresses(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userAddresses)
}

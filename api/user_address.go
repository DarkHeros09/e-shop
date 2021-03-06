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

type createUserAddressRequest struct {
	UserID      int64  `json:"user_id" binding:"required,min=1"`
	AddressLine string `json:"address_line" binding:"required"`
	City        string `json:"city" binding:"required"`
	Telephone   int32  `json:"telephone" binding:"required,numeric,min=910000000,max=929999999"`
}

func (server *Server) createUserAddress(ctx *gin.Context) {
	var req createUserAddressRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.CreateUserAddressParams{
		UserID:      authPayload.UserID,
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.GetUserAddressParams{
		ID:     req.ID,
		UserID: authPayload.UserID,
	}
	userAddress, err := server.store.GetUserAddress(ctx, arg)
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.ListUserAddressesParams{
		UserID: authPayload.UserID,
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

type updateUserAddressByUserIDRequest struct {
	ID          int64  `json:"id" binding:"required,min=1"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Telephone   int32  `json:"telephone" binding:"required,numeric,min=910000000,max=929999999"`
}

func (server *Server) updateUserAddressByUserID(ctx *gin.Context) {
	var req updateUserAddressByUserIDRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	arg := db.UpdateUserAddressByUserIDParams{
		ID:          req.ID,
		UserID:      authPayload.UserID,
		AddressLine: req.AddressLine,
		City:        req.City,
		Telephone:   req.Telephone,
	}

	userAddress, err := server.store.UpdateUserAddressByUserID(ctx, arg)
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

	ctx.JSON(http.StatusOK, userAddress)
}

type deleteUserAddressRequest struct {
	ID     int64 `json:"id" binding:"required,min=1"`
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

func (server *Server) deleteUserAddress(ctx *gin.Context) {
	var req deleteUserAddressRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.UserPayload)
	if req.UserID != authPayload.UserID {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err := server.store.DeleteUserAddress(ctx, req.ID)
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

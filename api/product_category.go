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

type createProductCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createProductCategoryRequest

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.AdminPayload)
	if authPayload.AdminID == 0 || authPayload.TypeID != 1 || !authPayload.Active {
		err := errors.New("account unauthorized")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateProductCategoryParams{
		Name:        req.Name,
		Description: req.Description,
	}

	Category, err := server.store.CreateProductCategory(ctx, arg)
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

	ctx.JSON(http.StatusOK, Category)
}

type getProductCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	var req getProductCategoryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	Category, err := server.store.GetProductCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, Category)
}

type listProductCategoriesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCategories(ctx *gin.Context) {
	var req listProductCategoriesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductCategoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	categories, err := server.store.ListProductCategories(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

type updateProductCategoryRequest struct {
	ID     int64 `json:"id" binding:"required,min=1"`
	Active *bool `json:"active" binding:"required"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var req updateProductCategoryRequest

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.AdminPayload)
	if authPayload.AdminID == 0 || authPayload.TypeID != 1 || !authPayload.Active {
		err := errors.New("account unauthorized")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateProductCategoryParams{
		ID:     req.ID,
		Active: *req.Active,
	}

	productCategory, err := server.store.UpdateProductCategory(ctx, arg)
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

	ctx.JSON(http.StatusOK, productCategory)
}

type deleteProductCategoryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var req deleteProductCategoryRequest

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.AdminPayload)
	if authPayload.AdminID == 0 || authPayload.TypeID != 1 || !authPayload.Active {
		err := errors.New("account unauthorized")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteProductCategory(ctx, req.ID)
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

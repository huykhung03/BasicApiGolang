package apii

import (
	"errors"
	"net/http"
	"simple_shop/db/sqlc"
	"simple_shop/token"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	ProductName   string `json:"product_name" binding:"required,min=1"`
	KindOfProduct string `json:"kind_of_product" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
	Price         int32  `json:"price" binding:"required,min=1"`
	Quantity      int32  `json:"quantity" binding:"required,min=1"`
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	admin, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if !admin.Level {
		err := errors.New("not admin")
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	product, err := server.store.CreateProduct(ctx,
		sqlc.CreateProductParams{
			ProductName:   req.ProductName,
			KindOfProduct: req.KindOfProduct,
			Owner:         authPayload.Username,
			Currency:      req.Currency,
			Price:         req.Price,
			Quantity:      req.Price,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

type deleteProductRequest struct {
	IDProduct int32 `json:"id_product" binding:"required,min=1"`
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	productsOfUser, err := server.store.ListProducts(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	var have bool = false

	for index := range productsOfUser {
		item := productsOfUser[index]
		if item.IDProduct == req.IDProduct {
			have = true
			break
		}
	}

	if have == false {
		err = errors.New("No product that you want to delete")
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	err = server.store.DeleteProduct(ctx, req.IDProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	message := "delete successfully"

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (server *Server) listProducts(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	productsOfUser, err := server.store.ListProducts(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, productsOfUser)
}

type getProductRequest struct {
	IDProduct int32 `uri:"id_product" binding:"required,min=1"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	productsOfUser, err := server.store.ListProducts(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	var have bool = false

	for index := range productsOfUser {
		item := productsOfUser[index]
		if item.IDProduct == req.IDProduct {
			have = true
			break
		}
	}

	if have == false {
		err = errors.New("No product that you want to get")
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, req.IDProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

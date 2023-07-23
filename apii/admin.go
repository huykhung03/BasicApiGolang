package apii

import (
	"net/http"
	"simple_shop/db/sqlc"
	"simple_shop/db/util"
	"simple_shop/token"

	"github.com/gin-gonic/gin"
)

func (server *Server) createAdmin(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	admin, err := server.store.CreateAdmin(ctx, sqlc.CreateAdminParams{
		Username:       req.Username,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		Level:          true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, admin)
}

func (server *Server) listPurchaseHistoriesOfAdmin(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	listProducts, err := server.store.ListProducts(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	listPurchaseHistories := []sqlc.PurchaseHistory{}

	for index := range listProducts {
		item := listProducts[index]
		purchaseHistories, err := server.store.GetPurchaseHistoriesByIdProduct(ctx, item.IDProduct)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		for index := range purchaseHistories {
			listPurchaseHistories = append(listPurchaseHistories, purchaseHistories[index])
		}
	}

	ctx.JSON(http.StatusOK, listPurchaseHistories)
}

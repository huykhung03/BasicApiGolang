package apii

import (
	"database/sql"
	"net/http"
	"simple_shop/db/sqlc"

	"github.com/gin-gonic/gin"
)

type purchaseRequest struct {
	IDProduct        int32  `json:"id_product" binding:"required, gte=1"`
	PurchaseQuantity int32  `json:"purchase_quantity" binding:"required, gte=1"`
	Buyer            string `json:"buyer" binding:"required"`
	CardNumber       string `json:"card_number" binding:"required"`
}

func (server *Server) createPurchase(ctx *gin.Context) {
	var req purchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	currency, err := server.store.GetCurrencyByCardNumberAndUsername(ctx, sqlc.GetCurrencyByCardNumberAndUsernameParams{
		Username:   req.Buyer,
		CardNumber: req.CardNumber,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if !server.validBankAccount(ctx, req.Buyer, currency) {
		return
	}

	product, err := server.store.GetProduct(ctx, req.IDProduct)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	result, err := server.store.PurchaseTransaction(ctx, sqlc.PurchaseTransactionPagrams{
		Product:           product,
		PurchaseQuantity:  uint16(req.PurchaseQuantity),
		Buyer:             req.Buyer,
		CardNumberOfBuyer: req.CardNumber,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

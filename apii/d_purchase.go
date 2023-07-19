package apii

import (
	"database/sql"
	"fmt"
	"net/http"
	"simple_shop/db/sqlc"

	"github.com/gin-gonic/gin"
)

type purchaseParams struct {
	IDProduct        int32  `json:"id_product" binding:"required"`
	PurchaseQuantity int32  `json:"purchase_quantity" binding:"required"`
	Buyer            string `json:"buyer" binding:"required"`
	CardNumber       string `json:"card_number" binding:"required"`
}

func (server *Server) validBankAccount(ctx *gin.Context, username string, currency string) (sqlc.BankAccount, bool) {
	bankAccount, err := server.store.GetBankAccountByUserNameAndCurrency(ctx,
		sqlc.GetBankAccountByUserNameAndCurrencyParams{
			Username: username,
			Currency: currency,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return bankAccount, false
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return bankAccount, false
	}

	if bankAccount.Currency != currency {
		err := fmt.Errorf("account [%s] currency mismatch: %s vs %s", bankAccount.Username, bankAccount.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return bankAccount, false
	}

	return bankAccount, true
}

func (server *Server) createPurchase(ctx *gin.Context) {
	var req purchaseParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
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

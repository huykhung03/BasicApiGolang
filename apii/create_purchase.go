package apii

import (
	"database/sql"
	"errors"
	"net/http"
	"simple_shop/db/sqlc"
	"simple_shop/token"

	"github.com/gin-gonic/gin"
)

// check checks if the product owner has a bank account that matches the currency of the product
func (server *Server) check(ctx *gin.Context, username string, currency string) (sqlc.BankAccount, bool) {
	bankAccount, err := server.store.GetBankAccountByUserNameAndCurrency(ctx,
		sqlc.GetBankAccountByUserNameAndCurrencyParams{
			Username: username,
			Currency: currency,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("The product owner does not have a bank account that matches the currency of the product")
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return bankAccount, false
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return bankAccount, false
	}

	return bankAccount, true
}

type purchaseRequest struct {
	IDProduct        int32  `json:"id_product" binding:"required,gte=1"`
	PurchaseQuantity int32  `json:"purchase_quantity" binding:"required,gte=1"`
	CardNumber       string `json:"card_number" binding:"required"`
}

func (server *Server) createPurchase(ctx *gin.Context) {
	var req purchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	product, err := server.store.GetProduct(ctx, req.IDProduct)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	_, valid := server.check(ctx, product.Owner, product.Currency)
	if !valid {
		return
	}

	result, err := server.store.PurchaseTransaction(ctx, sqlc.PurchaseTransactionPagrams{
		Product:           product,
		PurchaseQuantity:  uint16(req.PurchaseQuantity),
		Buyer:             authPayload.Username,
		CardNumberOfBuyer: req.CardNumber,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

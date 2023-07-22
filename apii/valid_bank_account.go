package apii

import (
	"database/sql"
	"fmt"
	"net/http"
	"simple_shop/db/sqlc"

	"github.com/gin-gonic/gin"
)

func (server *Server) validBankAccount(ctx *gin.Context, username string, currency string) bool {
	bankAccount, err := server.store.GetBankAccountByUserNameAndCurrency(ctx,
		sqlc.GetBankAccountByUserNameAndCurrencyParams{
			Username: username,
			Currency: currency,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	if bankAccount.Currency != currency {
		err := fmt.Errorf("Bank account [%s] currency mismatch: %s vs %s", bankAccount.Username, bankAccount.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}

	return true
}

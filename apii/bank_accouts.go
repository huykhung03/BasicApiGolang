package apii

import (
	"database/sql"
	"log"
	"net/http"
	"simple_shop/db/sqlc"
	"simple_shop/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createBankAccountResquest struct {
	Currency   string `json:"currency" binding:"required,currency"`
	CardNumber string `json:"card_number" binding:"required,min=8"`
	Balance    uint32 `json:"balance" binding:"required,min=1"`
}

func (server *Server) createBankAccount(ctx *gin.Context) {
	var req createBankAccountResquest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := sqlc.CreateBankAccountParams{
		Username:   authPayload.Username,
		Currency:   req.Currency,
		CardNumber: req.CardNumber,
		Balance:    int32(req.Balance),
	}

	bankAccount, err := server.store.CreateBankAccount(ctx, arg)
	if err != nil {
		if pqErr, oke := err.(*pq.Error); oke {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bankAccount)
}

type deleteBankAccountRequest struct {
	CardNumber string `json:"card_number" binding:"required"`
}

func (server *Server) deleteBankAccount(ctx *gin.Context) {
	var req deleteBankAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	_, err := server.store.GetCurrencyByCardNumberAndUsername(ctx,
		sqlc.GetCurrencyByCardNumberAndUsernameParams{
			CardNumber: req.CardNumber,
			Username:   authPayload.Username,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = server.store.DeleteBankAccountByCardNumberAndUserName(ctx,
		sqlc.DeleteBankAccountByCardNumberAndUserNameParams{
			CardNumber: req.CardNumber,
			Username:   authPayload.Username,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	message := "delete successfully"

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (server *Server) listBankAccounts(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	listBankAcounts, err := server.store.ListBankAccounts(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, listBankAcounts)
}

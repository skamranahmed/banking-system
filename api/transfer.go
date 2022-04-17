package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/skamranahmed/banking-system/db/sqlc"
	"github.com/skamranahmed/banking-system/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"` // the account from where the money is getting debited
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`   // the account to which the money is getting credited
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// extract the authPayload from the request context
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	// verify the currency of `fromAccount`
	fromAccount, isFromAccountValid := server.validAccount(c, req.FromAccountID, req.Currency)
	if !isFromAccountValid {
		return
	}

	// verify whether the `fromAccount` belongs to the authenticated user
	if fromAccount.UserID != int64(authPayload.UserID) {
		err := errors.New("fromAccount does not belong to the authenticated user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// verify whether the `fromAccount` has enough balance
	if fromAccount.Balance < req.Amount {
		err := fmt.Errorf("accountID: %d, insufficient balance: %d", fromAccount.ID, fromAccount.Balance)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify the currency of `toAccount`
	_, isToAccountValid := server.validAccount(c, req.ToAccountID, req.Currency)
	if !isToAccountValid {
		return
	}

	arg := db.TransferTxnParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTxn(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
	return
}

func (server *Server) validAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("no record found")))
			return account, false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("accountID:%d currency mismatch. Account Currency:%s, got currency:%s", account.ID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

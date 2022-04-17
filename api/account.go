package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/skamranahmed/banking-system/db/sqlc"
	"github.com/skamranahmed/banking-system/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// extract the authPayload from the request context
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		UserID:   int64(authPayload.UserID),
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(c, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, account)
	return
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// extract the authPayload from the request context
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := server.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("no record found")))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.UserID != uint(account.UserID) {
		err := errors.New("account does not belong to the authenticated user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
	return
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(c *gin.Context) {
	var req listAccountsRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// extract the authPayload from the request context
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		UserID: int64(authPayload.UserID),
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
	return
}

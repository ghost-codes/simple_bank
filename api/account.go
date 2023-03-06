package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/ghost-codes/simplebank/db/sqlc"
	"github.com/ghost-codes/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}
type GetAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type ListAccountRequest struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Count int32 `form:"count" binding:"required,min=5,max=20"`
}

func (server *Server) createAccount(context *gin.Context) {
	var req CreateAccountRequest
	err := context.BindJSON(&req)
	if err != nil {

		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(context, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "uniue_violation":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

func (server *Server) getAccountById(context *gin.Context) {
	var req GetAccountByIDRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(context, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			context.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != payload.Username {
		err := fmt.Errorf("account does not belong to authorized user: %v", payload.Username)
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(context *gin.Context) {
	var req GetAccountByIDRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.DeleteAccountParams{
		ID:    req.ID,
		Owner: payload.Username,
	}

	err := server.store.DeleteAccount(context, args)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			context.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	context.JSON(http.StatusOK, map[string]string{"message": "Account successfully deleted"})
}

func (server *Server) listAccounts(context *gin.Context) {
	var req ListAccountRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  req.Count,
		Offset: (req.Page - 1) * req.Count,
	}

	account, err := server.store.ListAccounts(context, args)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusOK, []int{})
		} else {
			context.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	context.JSON(http.StatusOK, account)
}

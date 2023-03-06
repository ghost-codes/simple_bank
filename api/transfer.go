package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/ghost-codes/simplebank/db/sqlc"
	"github.com/ghost-codes/simplebank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"fromAccountId" binding:"required,min=1"`
	ToAccountId   int64  `json:"toAccountId" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount, ok := server.validAccount(ctx, req.FromAccountId, req.Currency)

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != payload.Username {
		err := fmt.Errorf("account does not belong to authorized user: %v", payload.Username)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if !ok {
		return
	}

	if _, ok := server.validAccount(ctx, req.ToAccountId, req.Currency); !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (*db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return nil, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account %v currency mismatch:%v vs %v", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return nil, false
	}
	return &account, true
}

package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Support all queries and new functions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//execTx executes a funtion within a database transaction

func (store *Store) execTx(context context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(context, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v , rb err:%v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

// This struct contains all necessary information to transfer transactions into different accounts
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// this struct contains result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Transafer TX performs a money transfer from one account to another
// It creates a transfer record, add account entries
// And update account balance in a single database  transactions;
func (store *Store) TransferTx(context context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(context, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(context, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(context, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(context, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})

		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(context, q, args.FromAccountID, args.ToAccountID, -args.Amount, args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(context, q, args.ToAccountID, args.FromAccountID, args.Amount, -args.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	accountID2 int64,
	amount1 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	return

}

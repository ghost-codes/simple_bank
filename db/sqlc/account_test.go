package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ghost-codes/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	a := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), a.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, a, account)
}

func TestUpdateAccount(t *testing.T) {
	a := createRandomAccount(t)

	args := UpdateAcountParams{
		ID:      a.ID,
		Balance: util.RandomAmount(),
	}
	account, err := testQueries.UpdateAcount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.ID, args.ID)
	require.Equal(t, args.Balance, args.Balance)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)

}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

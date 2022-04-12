package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Account1 Balance Before Transaction(s): ", account1.Balance)
	fmt.Println("Account2 Balance Before Transaction(s): ", account2.Balance)
	fmt.Println("-----------------------------------------------------------")

	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxnResult)

	// run concurrent transfer transactions
	for i := 0; i < n; i++ {
		txnName := fmt.Sprintf("txn: %d", i+1)

		go func() {
			ctx := context.WithValue(context.Background(), txnKey, txnName)
			// transfer money from account1 to account2
			result, err := store.TransferTxn(ctx, TransferTxnParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check errors and results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account and account balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Printf("Account1 Balance After Transaction-%d: %d\n", i+1, fromAccount.Balance)
		fmt.Printf("Account2 Balance After Transaction-%d: %d\n", i+1, toAccount.Balance)

		// amount of money debited from account1 (initialAccountBalance - currentAccountBalance )
		debitedAmountAccount1 := account1.Balance - fromAccount.Balance

		// amount of money credited to account2 (currentAccountBalance - initalAccountBalance)
		creditedAmmountAccount2 := toAccount.Balance - account2.Balance

		// the amount of money going into account2 must be equal to the amount of money going out from account1
		require.Equal(t, creditedAmmountAccount2, debitedAmountAccount1)
		require.True(t, debitedAmountAccount1 > 0)
		require.True(t, debitedAmountAccount1%amount == 0)

		quotient := int(debitedAmountAccount1 / amount)
		require.True(t, quotient >= 1 && quotient <= n)
		require.NotContains(t, existed, quotient)
		existed[quotient] = true
	}

	// check the final updated balance of both the accounts
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Account1 Balance After Transaction(s): ", updatedAccount1.Balance)
	fmt.Println("Account2 Balance After Transaction(s): ", updatedAccount2.Balance)
	fmt.Println("-----------------------------------------------------------")

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

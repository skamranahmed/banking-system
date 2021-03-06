package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var txnKey = struct{}{}

// Store provides all functions to execute db queries and transaction
type Store interface {
	Querier
	TransferTxn(ctx context.Context, arg TransferTxnParams) (TransferTxnResult, error)
}

// SQLStore provides all functions to execute SQL queries and transaction
type SQLStore struct {
	*Queries
	db *sql.DB // this is required to create a new db txn
}

// NewStore : creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTxn : executes a function within a db transaction
func (s *SQLStore) execTxn(ctx context.Context, fn func(*Queries) error) error {
	// begin the transaction
	txn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(txn)
	err = fn(queries)
	if err != nil {
		// rollback the transaction
		rollbackErr := txn.Rollback()
		if rollbackErr != nil {
			errMsg := fmt.Sprintf("txn error: %v, rollback error: %v", err, rollbackErr)
			return errors.New(errMsg)
		}
		return err
	}

	return txn.Commit()
}

// TransferTxnParams : contains the input parameters of the transfer transaction
type TransferTxnParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxnResult : contains the result of transfer transaction
type TransferTxnResult struct {
	Transfer    Transfer `json:"transfer"`     // the created transfer record
	FromAccount Account  `json:"from_account"` // the `from account` after its balance has been updated
	ToAccount   Account  `json:"to_account"`   // the `to account` after its balance has been updated
	FromEntry   Entry    `json:"from_entry"`   // the entry record of the `from account`
	ToEntry     Entry    `json:"to_entry"`     // the entry record of the `to account`
}

// TransferTxn : performs money transfer from one account to the other
func (s *SQLStore) TransferTxn(ctx context.Context, arg TransferTxnParams) (TransferTxnResult, error) {
	/*
		Steps Involved:
		- Begin Transaction
			- Create a transfer record
			- Create individual entry records for both `from account` and `to account`
			- Update the balance of `from account`
			- Update the balance of `to account`
		- Commit
	*/

	var result TransferTxnResult

	err := s.execTxn(ctx, func(q *Queries) error {
		var err error
		txnName := ctx.Value(txnKey)

		fmt.Println(txnName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txnName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // debit
		})
		if err != nil {
			return err
		}

		fmt.Println(txnName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // credit
		})
		if err != nil {
			return err
		}

		// logic for updating the account balance of `from account` and `to account`
		/*
			// fmt.Println(txnName, "get account 1 for update")
			// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
			// if err != nil {
			// 	return err
			// }

			// fmt.Println(txnName, "update account 1")
			// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			// 	ID:      arg.FromAccountID,
			// 	Balance: account1.Balance - arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		*/

		// fmt.Println(txnName, "update account 1")
		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		/*
			// fmt.Println(txnName, "get account 2 for update")
			// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			// if err != nil {
			// 	return err
			// }

			// fmt.Println(txnName, "update account 2")
			// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			// 	ID:      arg.ToAccountID,
			// 	Balance: account2.Balance + arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		*/

		// fmt.Println(txnName, "update account 2")
		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		/*
			----------------------------------------------------------------------------------------------

			##################           Possible Deadlock Situation                ######################

			-- Transaction 1: transfer Rs.10 from account 1 to account 2
			BEGIN;

			UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
			UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;

			COMMIT;


			-- Transaction 2: transfer Rs.10 from account 2 to account 1
			BEGIN;

			UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;
			UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;

			COMMIT;

			----------------------------------------------------------------------------------------------

			##################           Solution for Preventing Deadlock              ###################



			if from_account_id < to_account_id -> We update the from_account first and then the to_account
			if the from_account_id > to_account_id -> We update the to_account first and then the from account

			-- Transaction 1: transfer Rs.10 from account 1 to account 2
			BEGIN;

			UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
			UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;

			COMMIT;


			-- Transaction 2: transfer Rs.10 from account 2 to account 1
			BEGIN;

			UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
			UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;

			COMMIT;

			----------------------------------------------------------------------------------------------

		*/

		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println(txnName, "updating the fromAccount")
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount, // debit
			})
			if err != nil {
				return err
			}

			fmt.Println(txnName, "updating the toAccount")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount, // credit
			})
			if err != nil {
				return err
			}

		} else {
			fmt.Println(txnName, "updating the toAccount")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount, // credit
			})
			if err != nil {
				return err
			}

			fmt.Println(txnName, "updating the fromAccount")
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount, // debit
			})
			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

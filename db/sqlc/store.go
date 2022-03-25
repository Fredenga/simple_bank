package db

import (
	"context"
	"database/sql"
	"fmt"
)

//provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}
func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

//execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID int32 `json:"to_account_id"`
	Amount int64  `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer 	`json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account 	`json:"to_account"`
	FromEntry Entry 	`json:"from_entry"`
	ToEntry Entry 		`json:"to_entry"`
}

var TxKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		TxName := ctx.Value(TxKey)

		fmt.Println(TxName, "Create Transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: int64(arg.FromAccountID),
			ToAccountID: int64(arg.ToAccountID),
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(TxName, "Create Entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: int64(arg.FromAccountID),
			Amount: arg.Amount,
		})
		if err != nil {return err}

		fmt.Println(TxName, "Create Entry 2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: int64(arg.ToAccountID),
			Amount: arg.Amount,
		})
		if err != nil {return err}
		
		fmt.Println(TxName, "Get Account1 for update")

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {return err}
			
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {return err}
		}


		return nil
	})
	return result, err

}
func AddMoney (ctx context.Context, q *Queries, accountID1 int32, amount1 int64, accountID2 int32, amount2 int64) (account1 Account, account2 Account, err error){
	account1, err = q.AddAccount(ctx, AddAccountParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return 
	}
	account2, err = q.AddAccount(ctx, AddAccountParams{
		ID: accountID2,
		Amount: amount2,
	})	
	if err != nil {
		return 
	}
	return

}
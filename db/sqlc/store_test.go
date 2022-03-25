package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTxDeadlock(t *testing.T){
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	// results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		TxName := fmt.Sprintf("tx %d", i+1)

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if(i % 2 == 1){
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func(){
			ctx := context.WithValue(context.Background(), TxKey, TxName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})

			errs <- err
			// results <- result
		}()
		err := <-errs
		require.NoError(t,err)

		// result := <-results 
		// require.NotEmpty(t, result)
		// existed := make(map[int]bool)


		// transfer := result.Transfer
		// require.NotEmpty(t, transfer)
		// require.Equal(t, account1.ID, int32(transfer.FromAccountID))
		// require.Equal(t, account2.ID, int32(transfer.ToAccountID))
		// require.Equal(t, -amount, -transfer.Amount)
		// require.NotZero(t, transfer.ID)
		// require.NotZero(t, transfer.CreatedAt)

		// _, err = store.GetTransfer(context.Background(), transfer.ID)
		// require.NoError(t,err)

		// fromEntry := result.FromEntry
		// require.NotEmpty(t, fromEntry)
		// require.Equal(t, account1.ID, int32(fromEntry.AccountID))
		// require.Equal(t, amount, fromEntry.Amount)
		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)

		// _, err = store.GetEntry(context.Background(), fromEntry.ID)
		// require.NoError(t,err)

		// toEntry := result.ToEntry
		// require.NotEmpty(t, toEntry)
		// require.Equal(t, account2.ID, int32(toEntry.AccountID))
		// require.Equal(t, amount, toEntry.Amount)
		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)

		// _, err = store.GetEntry(context.Background(), toEntry.ID)
		// require.NoError(t,err)

		// fromAccount := result.FromAccount
		// require.NotEmpty(t, fromAccount)
		// require.Equal(t, account1.ID, int32(fromAccount.ID))

		// toAccount := result.ToAccount
		// require.NotEmpty(t, toAccount)
		// require.Equal(t, account2.ID, int32(toAccount.ID))

		// fmt.Println(">>tx: ", fromAccount.Balance, toAccount.Balance)
		// diff1 := account1.Balance - fromAccount.Balance
		// diff2 := toAccount.Balance - account2.Balance
		// fmt.Println("after>> ", diff1, diff2)
		// require.Equal(t, diff1, diff2)
		// require.True(t, diff1 > 0)
		// require.True(t, diff1 % amount == 0)

		// k := int(diff1 / amount)
		// require.True(t, k >= 1 && k <= n)
		// require.NotContains(t, existed, k)
		// existed[k] = true
	}
	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), (account1.ID))
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), (account2.ID))
	require.NoError(t, err)

	fmt.Println(">>balance1", account1.Balance - int64(n) * amount )
	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)		
	require.Equal(t, account1.Balance , updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
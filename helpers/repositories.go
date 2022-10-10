package helpers

import (
	"log"
	"math/big"
	"strconv"

	"github.com/SoteriaTech/blockchain-functions/store"
)

// FindOrCreateBtcTransaction find a btc transaction and returns it, or create it if not exist and returns nothing
func FindOrCreateBtcTransaction(t *store.BtcTransactionSchema) (tx *store.BtcTransactionSchema, err error) {
	tx, err = store.Firestore.FindBtcTransaction(t.TxHash + strconv.Itoa(t.VoutIdx))
	if err != nil || tx != nil {
		return
	}

	err = store.Firestore.CreateBtcTransaction(t)
	return
}

// UpdateAccountBtcBalance update the btc balance of a user UID by a given amount
func UpdateAccountBtcBalance(uid string, amount *big.Float) (*big.Float, error) {
	bal, errBal := store.Firestore.FindBtcBalance(uid)
	if errBal != nil {
		return nil, errBal
	}
	newBalance := new(big.Float).Add(amount, big.NewFloat(bal))

	updatedBalance, errUpdate := store.Firestore.UpdateBtcBalance(uid, newBalance)
	if errUpdate != nil {
		return nil, errUpdate
	}

	return big.NewFloat(updatedBalance), nil
}

// ConfirmBtcTransactions confirm transactions and update corresponding balances
func ConfirmBtcTransactions(txs []*store.BtcTransactionSchema) (err error) {
	if err = store.Firestore.UpdateTransactionsConfirmation(txs); err != nil {
		log.Fatal(err)
		return
	}

	for _, t := range txs {
		a, err := store.Firestore.FindAccountByAddress(t.To)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if _, err = UpdateAccountBtcBalance(a.UID, big.NewFloat(t.Amount)); err != nil {
			log.Fatal(err)
			continue
		}
	}

	return
}

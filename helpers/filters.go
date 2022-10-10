package helpers

import (
	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/store"
)

// FilterTransactionsByAccountAddress filter a list of transactions by a list of btc addresses
func FilterTransactionsByAccountAddress(txs []*btc.Transaction, accs []*store.BtcAccountSchema) map[string]*store.BtcTransactionSchema {
	f := make(map[string]store.BtcAccountSchema, len(accs))
	out := make(map[string]*store.BtcTransactionSchema)
	for _, a := range accs {
		f[a.Address] = *a
	}
	for _, t := range txs {
		if acc, ok := f[t.Address]; ok {
			tx := &store.BtcTransactionSchema{
				To:          t.Address,
				TxHash:      t.Hash,
				Amount:      FromSatoshiToBtc(&t.Value),
				BlockHeight: t.BlockHeight,
				VoutIdx:     t.N,
			}

			out[acc.UID] = tx
		}
	}
	return out
}

// FilterTransactionsByHash filter transactions by a slice of hashes
func FilterTransactionsByHash(txs []*store.BtcTransactionSchema, hashes []string) (out []*store.BtcTransactionSchema) {
	f := make(map[string]*store.BtcTransactionSchema, len(txs))

	for _, t := range txs {
		f[t.TxHash] = t
	}

	for _, h := range hashes {
		if _, ok := f[h]; ok {
			out = append(out, f[h])
		}
	}
	return
}

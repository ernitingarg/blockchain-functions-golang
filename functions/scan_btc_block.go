package functions

import (
	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/helpers"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

// ScanBtcBlock scan a btc block for transactions
func ScanBtcBlock(height int, accs []*store.BtcAccountSchema) ([]*store.BtcAccountSchema, error) {

	// get transactions from 3 blocks earlier from store
	prevTxs, _ := store.Firestore.FindTransactionsFromBlockHeight(height - 3)
	if len(prevTxs) > 0 {
		var tbc []string
		for _, t := range prevTxs {
			tbc = append(tbc, t.TxHash)
		}
		hashes, _ := btc.BtcService.ConfirmTransactions(tbc)
		if len(hashes) > 0 {
			cTxs := helpers.FilterTransactionsByHash(prevTxs, hashes)
			if err := helpers.ConfirmBtcTransactions(cTxs); err != nil {
				utils.ErrorReport.LogAndPrintError(err)
			}
		}
	}

	txs, err := btc.BtcService.ScanBlock(height)
	if err != nil {
		return nil, err
	}

	walletTxs := helpers.FilterTransactionsByAccountAddress(txs, accs)
	var uaccs []*store.BtcAccountSchema
	for uid, t := range walletTxs {
		t.Confirmed = false
		exists, errTx := helpers.FindOrCreateBtcTransaction(t)
		if errTx != nil {
			return nil, errTx
		}
		if exists != nil {
			continue
		}
		uaccs = append(uaccs, &store.BtcAccountSchema{UID: uid, Address: t.To, Balance: t.Amount})
	}

	return uaccs, nil
}

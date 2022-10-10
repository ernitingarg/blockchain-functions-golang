package functions

import (
	"math/big"

	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/helpers"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

// SyncBtcBalance sync the balance of user's account from its uid
func SyncBtcBalance(uid string) (*store.BtcAccountSchema, *utils.ErrorService) {
	btcAccount, errFind := store.Firestore.FindBtcAccount(uid)
	if errFind != nil {
		return nil, &utils.ErrorService{Code: 404, Err: errFind}
	}
	newBalance, errBalance := btc.BtcService.GetAccountBalance(btcAccount.Address)
	if errBalance != nil {
		return nil, &utils.ErrorService{Code: 400, Err: errBalance}
	}

	updatedBalance, errUpdate := store.Firestore.UpdateBtcBalance(btcAccount.UID, big.NewFloat(helpers.FromSatoshiToBtc(newBalance)))
	if errUpdate != nil {
		return nil, &utils.ErrorService{Code: 400, Err: errUpdate}
	}

	btcAccount.Balance = updatedBalance
	return btcAccount, nil
}

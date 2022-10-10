package store

import "time"

//BtcAccountSchema firestore schema of a bitcoin account
type BtcAccountSchema struct {
	UID     string  `json:"uid"`
	Address string  `json:"address"`
	Balance float64 `json:"BTC"`
}

// BtcTransactionSchema firestore schema of a btc transaction
type BtcTransactionSchema struct {
	Amount      float64 `firestore:"amount"`
	To          string  `firestore:"to"`
	TxHash      string  `firestore:"txHash"`
	VoutIdx     int     `firestore:"vout_idx"`
	BlockHeight int     `firestore:"block_height"`
	Confirmed   bool    `firestore:"confirmed"`
}

// ChainStateSchema firestore schema of a chain state
type ChainStateSchema struct {
	Hash        string    `firestore:"hash"`
	Time        int       `firestore:"time"`
	LastUpdated time.Time `firestore:"last_updated"`
	BlockIndex  int       `firestore:"block_index"`
	Height      int       `firestore:"height"`
	TxIndexes   []int     `firestore:"txIndexes"`
}

package btc

import (
	"math/big"
	"time"
)

// HeadBlock structure of the head block
type HeadBlock struct {
	Hash        string    `firestore:"hash"`
	Time        int       `firestore:"time"`
	LastUpdated time.Time `firestore:"last_updated"`
	BlockIndex  int       `firestore:"block_index"`
	Height      int       `firestore:"height"`
	TxIndexes   []int     `firestore:"txIndexes"`
}

// Block structure of a BTC block
type Block struct {
	Hash         string `json:"hash"`
	Ver          int    `json:"ver"`
	PrevBlock    string `json:"prev_block"`
	MrklRoot     string `json:"mrkl_root"`
	Time         int    `json:"time"`
	Nonce        int    `json:"nonce"`
	NTx          int    `json:"n_tx"`
	BlockIndex   int    `json:"block_index"`
	MainChain    bool   `json:"main_chain"`
	Height       int    `json:"height"`
	ReceivedTime int    `json:"received_time"`
	Txs          []Tx   `json:"tx"`
}

// Address structure address of a BTC account
type Address struct {
	Hash160       string `json:"hash160"`
	Address       string `json:"address"`
	NTx           int    `json:"n_tx"`
	TotalReceived int    `json:"total_received"`
	TotalSent     int    `json:"total_sent"`
	FinalBalance  int    `json:"final_balance"`
	Txs           []*Tx  `json:"txs"`
}

// Transaction decoded transaction from TX inputs and outputs with only required properties
type Transaction struct {
	Address     string  `json:"address"`
	Value       big.Int `json:"value"`
	BlockHeight int     `json:"block_height"`
	Hash        string  `json:"hash"`
	TxIndex     big.Int `json:"tx_index"`
	N           int     `json:"n"`
}

// Tx structure of a BTC transaction
type Tx struct {
	Result      int       `json:"result"`
	Ver         int       `json:"ver"`
	Size        int       `json:"size"`
	Inputs      []*Inputs `json:"inputs"`
	Time        int       `json:"time"`
	BlockHeight int       `json:"block_height"`
	Fee         big.Int   `json:"fee"`
	TxIndex     int       `json:"tx_index"`
	VinSz       int       `json:"vin_sz"`
	Hash        string    `json:"hash"`
	VoutSz      int       `json:"vout_sz"`
	RelayedBy   string    `json:"relayed_by"`
	Out         []*Out    `json:"out"`
}

// Inputs inputs of a BTC transaction
type Inputs struct {
	Sequence int     `json:"sequence"`
	Script   string  `json:"script"`
	PrevOut  PrevOut `json:"prev_out"`
}

// PrevOut PrevOut of a btc Input
type PrevOut struct {
	Spent   bool    `json:"spent"`
	TxIndex big.Int `json:"tx_index"`
	Type    int     `json:"type"`
	Addr    string  `json:"addr"`
	Value   big.Int `json:"value"`
	N       int     `json:"n"`
	Script  string  `json:"script"`
}

// Out Out of a btc transaction
type Out struct {
	Spent   bool    `json:"spent"`
	TxIndex big.Int `json:"tx_index"`
	Type    int     `json:"type"`
	Addr    string  `json:"addr"`
	Value   big.Int `json:"value"`
	N       int     `json:"n"`
	Script  string  `json:"script"`
}

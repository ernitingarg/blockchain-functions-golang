package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/utils"

	"github.com/blockcypher/gobcy"
)

const (
	token   string = "8dec31da4fbc4b6794ca4abc8eb77bf1"
	testnet        = "test3"
	mainnet        = "main"
)

// BlockCypherClient structure of the blockCypher client
type BlockCypherClient struct {
	client gobcy.API
	http   *http.Client
}

// BlockCypher BlockCypher client instance
var BlockCypher *BlockCypherClient

// InitBlockCypherClient initialize an instance of BlockCypher. Chain is either "main" or "test3"
func InitBlockCypherClient() {
	client := gobcy.API{Token: token, Coin: "btc", Chain: getChain()}

	BlockCypher = &BlockCypherClient{
		client: client,
		http:   &http.Client{},
	}
}

//GetBalance get the balance of a given account
func (b *BlockCypherClient) GetBalance(address string) (*big.Float, error) {

	acc, err := b.client.GetAddrBal(address, nil)
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		return nil, err
	}
	balance := acc.Balance
	return new(big.Float).SetInt(&balance), nil
}

// FetchTransactionsFromBlock fetch
func (b *BlockCypherClient) FetchTransactionsFromBlock(height int) ([]*gobcy.TX, error) {
	block, err := b.GetBlock(height)
	if err != nil {
		return nil, err
	}

	txs, errs := b.aggregateTransactions(block.TXids)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return txs, nil
}

// GetBlock get the data of a given block, if provided height is 0 then get the head block
func (b *BlockCypherClient) GetBlock(height int) (*gobcy.Block, error) {
	var block gobcy.Block
	var err error

	if height == 0 {
		height, err = b.getHeadBlock()
	}

	block, err = b.client.GetBlock(height, "", map[string]string{"limit": "500"})
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		return nil, err
	}

	return &block, nil
}

//AggregateTransactions aggregate transactions from a list of txIds
func (b *BlockCypherClient) aggregateTransactions(txIDs []string) ([]*gobcy.TX, []error) {
	var txs []*gobcy.TX
	var errs []error

	for _, txID := range txIDs {
		time.Sleep(2 * time.Second)
		tx, err := b.GetTransaction(txID)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		txs = append(txs, tx)
	}
	return txs, errs
}

// GetTransaction get details of the transaction with the given txId
func (b *BlockCypherClient) GetTransaction(txID string) (*gobcy.TX, error) {
	tx := &gobcy.TX{}
	if err := b.RequestTxURL(txID, tx); err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		return nil, err
	}
	return tx, nil
}

// RequestTxURL make a request to the transaction URL for transaction details
func (b *BlockCypherClient) RequestTxURL(txID string, i interface{}) error {
	rsp, err := b.http.Get(getBaseTxURL() + txID + "?token=" + token)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if rsp.Status[0] != '2' {
		return fmt.Errorf("expected status 2xx, got %s: %s", rsp.Status, string(data))
	}

	return json.Unmarshal(data, &i)
}

func (b *BlockCypherClient) getHeadBlock() (int, error) {
	chain, err := b.client.GetChain()
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		return 0, err
	}
	return chain.Height, nil
}

func getBaseTxURL() string {
	if env.EnvVars.ProjectID != env.PRODUCTION {
		return "https://api.blockcypher.com/v1/btc/test3/txs/"
	}
	return "https://api.blockcypher.com/v1/btc/main/txs/"
}

func getChain() string {
	if env.EnvVars.ProjectID == env.PRODUCTION {
		return mainnet
	}
	return testnet
}

func formatBlock(b *gobcy.Block) *btc.Block {
	return &btc.Block{
		Hash:         b.Hash,
		Height:       b.Height,
		Ver:          b.Ver,
		PrevBlock:    b.PrevBlock,
		MrklRoot:     b.MerkleRoot,
		NTx:          b.NumTX,
		Nonce:        b.Nonce,
		Time:         int(b.Time.Unix()),
		ReceivedTime: int(b.ReceivedTime.Unix()),
	}
}

func formatTx(t *gobcy.TX) *btc.Tx {
	return &btc.Tx{
		Hash:        t.Hash,
		BlockHeight: t.BlockHeight,
		Time:        int(t.Received.Unix()),
		Fee:         t.Fees,
		VinSz:       t.VinSize,
		VoutSz:      t.VoutSize,
	}
}

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"

	"github.com/SoteriaTech/blockchain-functions/btc"
)

const (
	baseURL string = "https://blockchain.info"
)

// BlockInfoClient structure of the blockInfo api client
type BlockInfoClient struct {
	*http.Client
}

type bIAccount struct {
	Address       string  `json:"address"`
	FinalBalance  big.Int `json:"final_balance"`
	TotalReceived big.Int `json:"total_received"`
	TotalSent     big.Int `json:"total_sent"`
	NTx           big.Int `json:"n_tx"`
	NUnredeemed   big.Int `json:"n_unredeemed"`
}

// BlockInfo instance of the BlockInfoClient api
var BlockInfo *BlockInfoClient

// InitBlockInfoClient initialize an instance of BlockInfo
func InitBlockInfoClient() {
	BlockInfo = &BlockInfoClient{
		Client: &http.Client{},
	}
}

// GetBalance get the balance of the account corresponding to the given address
func (b *BlockInfoClient) GetBalance(address string) (*big.Int, error) {
	acc := &bIAccount{}
	endpoint := "/rawaddr/" + address
	err := b.request(endpoint, acc, true)
	if err != nil {
		return nil, err
	}

	return &acc.FinalBalance, nil
}

// GetHeadBlock get the head block basic info
func (b *BlockInfoClient) GetHeadBlock() (*btc.HeadBlock, error) {
	lb := &btc.HeadBlock{}
	err := b.request("/latestblock", lb, true)
	if err != nil {
		return nil, err
	}

	return lb, nil
}

// GetBlock get block fat given height or, if height is 0, get head block
func (b *BlockInfoClient) GetBlock(height int) (*btc.Block, error) {
	block := &btc.Block{}
	endpoint := "/rawblock/" + strconv.Itoa(height)
	err := b.request(endpoint, block, true)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetTransactionsFromBlock extract and parse transactions from a given block
func (b *BlockInfoClient) GetTransactionsFromBlock(block *btc.Block) ([]*btc.Transaction, []error) {
	var txs []*btc.Transaction
	var errs []error
	for _, tx := range block.Txs {
		txs = append(txs, parseTx(&tx, block.Height)...)
	}
	return txs, errs
}

// GetTransactionByHash get the detail of a transaction from its hash
func (b *BlockInfoClient) GetTransactionByHash(hash string) (tx *btc.Transaction, err error) {
	endpoint := "/rawtx/" + hash

	if err = b.request(endpoint, &tx, true); err != nil {
		return nil, err
	}
	return tx, nil
}

func (b *BlockInfoClient) request(endpoint string, i interface{}, isJSON bool) error {
	fullPath := baseURL + endpoint
	if isJSON {
		fullPath = baseURL + endpoint + "?format=json"
	}

	rsp, err := b.Get(fullPath)
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

func parseTx(tx *btc.Tx, height int) (ts []*btc.Transaction) {

	outs := parseOutTxs(tx.Out, tx.Hash, height)
	ts = append(ts, outs...)

	return
}

func parseOutTxs(out []*btc.Out, hash string, height int) (ts []*btc.Transaction) {
	for _, o := range out {
		t := &btc.Transaction{
			Hash:        hash,
			Address:     o.Addr,
			Value:       o.Value,
			TxIndex:     o.TxIndex,
			N:           o.N,
			BlockHeight: height,
		}
		ts = append(ts, t)
	}

	return
}

func parseInTxs(in []*btc.Inputs, hash string, height int) (ts []*btc.Transaction) {
	for _, i := range in {
		if i.PrevOut.Value.BitLen() == 0 {
			continue
		}
		sentValue := new(big.Int).Neg(&i.PrevOut.Value)
		t := &btc.Transaction{
			Hash:        hash,
			Address:     i.PrevOut.Addr,
			Value:       *sentValue,
			TxIndex:     i.PrevOut.TxIndex,
			N:           i.PrevOut.N,
			BlockHeight: height,
		}
		ts = append(ts, t)
	}

	return
}

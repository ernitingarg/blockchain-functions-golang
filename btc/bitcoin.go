package btc

import (
	"math/big"

	"github.com/blockcypher/gobcy"
)

type result struct {
	string
	TX *gobcy.TX
	error
}

// BitcoinAPI interface that the Btc Service implements
type BitcoinAPI interface {
	GetBlock(height int) (*Block, error)
	GetHeadBlock() (*HeadBlock, error)
	GetTransactionsFromBlock(block *Block) ([]*Transaction, []error)
	GetTransactionByHash(hash string) (*Transaction, error)
	GetBalance(address string) (*big.Int, error)
}

//Btc structure of the Btc service
type Btc struct {
	api BitcoinAPI
}

// BtcService instance of the btc service
var BtcService *Btc

// InitBtcService initialize the instance of the btc service
func InitBtcService(a BitcoinAPI) {
	BtcService = &Btc{
		api: a,
	}
}

// FetchBlock fetch block with given height. If height is 0, then fetch head block
func (b *Btc) FetchBlock(height int) (*Block, error) {

	block, err := b.api.GetBlock(height)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetAccountBalance get the balance of the account corresponding to the given address
func (b *Btc) GetAccountBalance(address string) (*big.Int, error) {
	balance, err := b.api.GetBalance(address)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// ScanBlock scan a btc Block, extract and parse its transactions
func (b *Btc) ScanBlock(height int) ([]*Transaction, error) {
	block, err := b.FetchBlock(height)
	if err != nil {
		return nil, err
	}
	txs, errs := b.api.GetTransactionsFromBlock(block)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return txs, nil
}

// GetHeadInfo get the info of the head block of the blockchain
func (b *Btc) GetHeadInfo() (*HeadBlock, error) {
	lb, err := b.api.GetHeadBlock()

	if err != nil {
		return nil, err
	}

	return lb, err
}

// ConfirmTransactions ask the blockchain for confirmed transactions
func (b *Btc) ConfirmTransactions(hashes []string) (confirmed []string, errs []error) {
	for _, h := range hashes {
		tx, err := b.api.GetTransactionByHash(h)

		if err != nil {
			errs = append(errs, err)
			continue
		}
		confirmed = append(confirmed, tx.Hash)
	}
	return
}

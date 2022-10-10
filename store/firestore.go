package store

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/env"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// FireStoreStore struct for firestore DB
type FireStoreStore struct {
	Client *firestore.Client
	ctx    context.Context
}

// Firestore instance of Firestore store
var Firestore *FireStoreStore

//InitFirestoreStore initialize a new firestore client
func InitFirestoreStore() {
	ctx := context.Background()
	client := newFireStoreClient(ctx)

	Firestore = &FireStoreStore{
		Client: client,
		ctx:    ctx,
	}
}

func newFireStoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, env.EnvVars.ProjectID, option.WithCredentialsFile(env.EnvVars.Keypath))
	if err != nil {
		log.Fatalf("Failed to create firestore client %v", err)
	}
	return client
}

// FindBtcAccount find btc account from a user UID
func (f *FireStoreStore) FindBtcAccount(uid string) (*BtcAccountSchema, error) {
	var btcAccount *BtcAccountSchema

	doc, err := f.Client.Collection("btc_accounts").Doc(uid).Get(f.ctx)
	if err != nil {
		return nil, err
	}
	if err := doc.DataTo(&btcAccount); err != nil {
		return nil, err
	}
	btcAccount.UID = uid

	return btcAccount, nil
}

// FindBtcBalance find the btc balance of a user UID
func (f *FireStoreStore) FindBtcBalance(uid string) (float64, error) {
	data := make(map[string]float64)

	doc, err := f.Client.Collection("balances").Doc(uid).Get(f.ctx)
	if err != nil {
		return 0, err
	}

	if err := doc.DataTo(&data); err != nil {
		return 0, err
	}

	return data["BTC"], nil
}

// UpdateBtcBalance update the btc balance of a user's account
func (f *FireStoreStore) UpdateBtcBalance(uid string, newBalance *big.Float) (float64, error) {
	doc := make(map[string]interface{})
	flBalance, _ := newBalance.Float64()
	doc["BTC"] = flBalance

	_, err := f.Client.Collection("balances").Doc(uid).Set(f.ctx, doc, firestore.MergeAll)
	if err != nil {
		return flBalance, err
	}
	return flBalance, nil
}

// GetAllAccountAddresses get all the current accounts and addresses from Soteria
func (f *FireStoreStore) GetAllAccountAddresses() ([]*BtcAccountSchema, error) {
	var accs []*BtcAccountSchema
	iter := f.Client.Collection("btc_accounts").Documents(f.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var acc *BtcAccountSchema
		doc.DataTo(&acc)
		acc.UID = doc.Ref.ID

		accs = append(accs, acc)
	}
	return accs, nil
}

// GetConvertRequests get all convert requests of a given account
func (f *FireStoreStore) GetConvertRequests(uid string) (map[string]interface{}, error) {
	docs := make(map[string]interface{})
	iter := f.Client.Collection("convert_history").Doc(uid).Collection("history").Documents(f.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		docs[uid] = doc.Data()
	}

	return docs, nil
}

// FindBtcTransaction find a btc transaction by hash
func (f *FireStoreStore) FindBtcTransaction(idx string) (acc *BtcTransactionSchema, err error) {
	doc, errStore := f.Client.Collection("btc_transactions").Doc(idx).Get(f.ctx)

	if errStore != nil && grpc.Code(errStore) != codes.NotFound {
		err = errStore
		return
	}

	if doc.Exists() {
		errDoc := doc.DataTo(&acc)
		if errDoc != nil {
			err = errDoc
		}
	}
	return
}

// CreateBtcTransaction create a btc transaction
func (f *FireStoreStore) CreateBtcTransaction(t *BtcTransactionSchema) (err error) {
	_, err = f.Client.Collection("btc_transactions").Doc(t.TxHash+strconv.Itoa(t.VoutIdx)).Set(f.ctx, &t)
	return
}

// GetChainState get the latest block data of the given chain from the store
func (f *FireStoreStore) GetChainState(chain string) (*btc.HeadBlock, error) {
	var hs *btc.HeadBlock
	doc, err := f.Client.Collection("chain_state").Doc(chain).Get(f.ctx)
	if err != nil {
		return nil, err
	}
	errData := doc.DataTo(&hs)
	if errData != nil {
		return nil, errData
	}
	return hs, err
}

// UpdateChainState update the latest block data of the given chain
func (f *FireStoreStore) UpdateChainState(chain string, data *btc.HeadBlock) (err error) {
	doc := make(map[string]interface{})
	doc["height"] = data.Height
	doc["time"] = data.Time
	doc["last_updated"] = time.Now()
	doc["block_index"] = data.BlockIndex
	doc["tx_indexes"] = data.TxIndexes

	_, err = f.Client.Collection("chain_state").Doc(chain).Set(f.ctx, doc)
	return
}

// FindTransactionsFromBlockHeight find transactions that have been recorded from a specific block height
func (f *FireStoreStore) FindTransactionsFromBlockHeight(h int) (txs []*BtcTransactionSchema, err error) {
	iter := f.Client.Collection("btc_transactions").Where("block_height", "==", h).Where("confirmed", "==", false).Documents(f.ctx)
	for {
		doc, errIter := iter.Next()
		if errIter == iterator.Done {
			break
		}
		if errIter != nil {
			err = errIter
			return
		}
		var tx *BtcTransactionSchema
		if err = doc.DataTo(&tx); err != nil {
			return
		}
		txs = append(txs, tx)
	}

	return
}

// UpdateTransactionsConfirmation update confirmation for each given transaction
func (f *FireStoreStore) UpdateTransactionsConfirmation(txs []*BtcTransactionSchema) (err error) {
	for _, t := range txs {
		uid := t.TxHash
		if t.VoutIdx >= 0 {
			uid = uid + strconv.Itoa(t.VoutIdx)
		}
		_, errSet := f.Client.Collection("btc_transactions").Doc(uid).Set(f.ctx, BtcTransactionSchema{Confirmed: true}, firestore.Merge([]string{"confirmed"}))
		if err != nil {
			err = errSet
			continue
		}
	}

	return
}

// FindAccountByAddress find a firestore account from an address
func (f *FireStoreStore) FindAccountByAddress(addr string) (a *BtcAccountSchema, err error) {
	doc, errQ := f.Client.Collection("btc_accounts").Where("address", "==", addr).Documents(f.ctx).Next()
	if errQ != nil {
		err = errQ
		return
	}
	if doc.Exists() {
		doc.DataTo(&a)
		a.UID = doc.Ref.ID
	}
	return
}

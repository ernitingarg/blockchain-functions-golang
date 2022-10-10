package functions

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/SoteriaTech/blockchain-functions/api"
	"github.com/SoteriaTech/blockchain-functions/btc"
	"github.com/SoteriaTech/blockchain-functions/env"
	"github.com/SoteriaTech/blockchain-functions/functions"
	"github.com/SoteriaTech/blockchain-functions/store"
	"github.com/SoteriaTech/blockchain-functions/utils"
)

const jsonContentType = "application/json"

// init function is ran automatically by GCP prior to the rest
func init() {
	env.InitEnvVars()
	utils.InitErrorReporting(env.EnvVars.ProjectID)
	store.InitFirestoreStore()
	api.InitBlockInfoClient()
	btc.InitBtcService(api.BlockInfo)
}

/***********************************************
*
* HTTP functions
*
***********************************************/

// SyncBtcBalance function sync the btc balance of a given user's account
func SyncBtcBalance(w http.ResponseWriter, r *http.Request) {

	data, errReq := utils.RequestData(r)
	if errReq != nil {
		utils.ErrorReport.LogAndPrintError(errReq)
		utils.RespondJSONWithError(w, 400, errReq.Error())
	}

	btcAccount, err := functions.SyncBtcBalance(data["uid"])
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err.Err)
		utils.RespondJSONWithError(w, err.Code, err.Err.Error())
	}
	utils.RespondJSON(w, 200, btcAccount)
}

// ScanBtcBlock scan a bitcoin blockchain block and parse it
func ScanBtcBlock(w http.ResponseWriter, r *http.Request) {
	data, errReq := utils.RequestData(r)
	if errReq != nil {
		utils.RespondJSONWithError(w, 400, errReq.Error())
	}

	height, errConv := strconv.Atoi(data["height"])
	if errConv != nil {
		utils.ErrorReport.LogAndPrintError(errConv)
		utils.RespondJSONWithError(w, 400, "error height format is incorrect")
		return
	}

	accs, errAccs := store.Firestore.GetAllAccountAddresses()
	if errAccs != nil {
		utils.RespondJSONWithError(w, 500, errAccs.Error())
		return
	}

	rsp, err := functions.ScanBtcBlock(height, accs)
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		utils.RespondJSONWithError(w, 500, err.Error())
	}

	utils.RespondJSON(w, 200, rsp)
}

// ScanBtcHead scan this is a replica of the pub/sub to test on the local server
func ScanBtcHead(w http.ResponseWriter, r *http.Request) {
	chain := env.EnvVars.BtcChain
	cs, err := store.Firestore.GetChainState(chain)
	if err != nil {
		utils.RespondJSONWithError(w, 500, err.Error())
		return
	}

	headBlock, err := btc.BtcService.GetHeadInfo()
	if err != nil {
		utils.RespondJSONWithError(w, 500, err.Error())
		return
	}

	if cs.Height == headBlock.Height {
		utils.RespondJSON(w, 200, cs)
		return
	}
	accs, errAccs := store.Firestore.GetAllAccountAddresses()
	if errAccs != nil {
		utils.RespondJSONWithError(w, 500, errAccs.Error())
		return
	}
	currHeight := cs.Height
	var blocks []int
	for {
		currHeight++

		_, errScan := functions.ScanBtcBlock(currHeight, accs)
		if errScan != nil {
			utils.ErrorReport.LogAndPrintError(errScan)
			utils.RespondJSONWithError(w, 500, errScan.Error())
		}
		blocks = append(blocks, currHeight)
		if currHeight == headBlock.Height {
			break
		}
	}

	errUpdate := store.Firestore.UpdateChainState(chain, headBlock)
	if errUpdate != nil {
		utils.RespondJSONWithError(w, 500, errUpdate.Error())
		return
	}

	utils.RespondJSON(w, 200, blocks)
	return
}

/***********************************************
*
* Pub/Sub functions
*
***********************************************/

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ScanBtcPubSub ping the btc blockchain for new block and scan them for transactions
func ScanBtcPubSub(ctx context.Context, m PubSubMessage) error {

	chain := env.EnvVars.BtcChain
	cs, err := store.Firestore.GetChainState(chain)
	if err != nil {
		return err
	}

	headBlock, err := btc.BtcService.GetHeadInfo()
	if err != nil {
		utils.ErrorReport.LogAndPrintError(err)
		return err
	}

	// if the head block hasn't changed we do nothing
	if cs.Height == headBlock.Height {
		return nil
	}

	//TODO : handle case if headBlock.Height < cs.Height (eg. reorg), if that's the case we may want to rollback some stuff

	accs, errAccs := store.Firestore.GetAllAccountAddresses()
	if errAccs != nil {
		utils.ErrorReport.LogAndPrintError(errAccs)
		return errAccs
	}
	currHeight := cs.Height
	var blocks []int

	//  loop through every block missing between our last state and the blockchain state
	for {
		currHeight++

		_, errScan := functions.ScanBtcBlock(currHeight, accs)
		if errScan != nil {
			utils.ErrorReport.LogAndPrintError(errScan)
			// we stop right here if we get an error
			return errScan
		}

		blocks = append(blocks, currHeight)
		if currHeight == headBlock.Height {
			break
		}
	}

	errUpdate := store.Firestore.UpdateChainState(chain, headBlock)
	if errUpdate != nil {
		utils.ErrorReport.LogAndPrintError(errUpdate)
		return errUpdate
	}
	log.Printf("Blocks  aggregated: %v", blocks)
	return nil
}

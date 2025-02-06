package guardian

import (
	"encoding/json"

	"github.com/catalogfi/blockchain/btc"
)

type Batch struct {
	Tx btc.Transaction
	// maps RequestId to Vouts in Tx
	RequestIds       map[string]int
	PreviousBatchID  string
	// maps invalidateTxId to the mergeRequestID that invalidates it
	MergeMap	map[string]string
	MergeTxFee       int64
	LastFailedTxHash string
}

const CoinbaseBatchID = "coinbase"
const GuardianWitnessSize = 267
const GuardianChangeSize = 43

func NewBatch(tx btc.Transaction, requestIds map[string]int, previousBatchID string, mergeTxFee int64) *Batch {
	return &Batch{
		Tx:              tx,
		RequestIds:      requestIds,
		PreviousBatchID: previousBatchID,
		MergeTxFee:      mergeTxFee,
	}
}

// implement marshalling and unmarshalling
func (b *Batch) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Batch) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}

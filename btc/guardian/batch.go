package guardian

import (
	"encoding/json"

	"github.com/catalogfi/blockchain/btc"
)

type Batch struct {
	Tx              btc.Transaction
	RequestIds      []string
	PreviousBatchID string
}

func NewBatch(tx btc.Transaction, requestIds []string, previousBatchID string) *Batch {
	return &Batch{
		Tx:              tx,
		RequestIds:      requestIds,
		PreviousBatchID: previousBatchID,
	}
}

// implement marshalling and unmarshalling
func (b *Batch) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Batch) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}

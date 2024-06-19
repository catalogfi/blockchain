package btc

type TxInputRequest struct {
	Witness   []byte
	Utxo      UTXO
	ToAddress string
	isSACP    bool
}

type TxOutputRequest struct {
	Value   uint64
	Address string
}

type Wallet interface {
	AddInput(tx *TxInputRequest) (string, error)
	AddOutput(tx TxOutputRequest) (string, error)
	Status(id string) (Transaction, bool, error)
}

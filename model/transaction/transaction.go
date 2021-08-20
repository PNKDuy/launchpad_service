package transaction

type Transaction struct {
	Hash string `json:"hash"`
	Version int `json:"ver,omitempty"`
	VInputSize int `json:"vin_sz,omitempty"`
	VOutputSize int `json:"vout_sz,omitempty"`
	Size int `json:"size,omitempty,omitempty"`
	Weight int `json:"weight,omitempty"`
	Fee int `json:"fee,omitempty"`
	RelayedBy string `json:"relayed_by,omitempty"`
	LockTime int64 `json:"lock_time,omitempty"`
	TxIndex int64 `json:"tx_index,omitempty"`
	DoubleSpend bool `json:"double_spend,omitempty"`
	Time int64 `json:"time,omitempty"`
	BlockIndex int `json:"block_index,omitempty"`
	BlockHeight int `json:"block_height,omitempty"`
	Inputs []Input `json:"inputs,omitempty"`
	Output []Output `json:"out,omitempty"`
}

type Input struct {
	Sequence int64 `json:"sequence"`
	Witness string `json:"witness"`
	Script string `json:"script"`
	Index int `json:"index"`
	PrevOut Output `json:"prev_out"`
}

type Output struct {
	Type int `json:"type,omitempty"`
	Spent bool `json:"spent,omitempty"`
	Value int64 `json:"value,omitempty"`
	SpendingOutpoints []*Output `json:"spending_outpoints,omitempty"`
	N int `json:"n"`
	TxIndex int64 `json:"tx_index"`
	Script string `json:"script,omitempty"`
	Address string `json:"addr,omitempty"`
}
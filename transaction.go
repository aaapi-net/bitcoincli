package bitcoincli

type Transaction struct {
	Amount float64               `json:"amount"`
	Fee float64                  `json:"fee"`
	Confirmations int            `json:"confirmations"`
	BlockHash string             `json:"blockhash"`
	BlockIndex int               `json:"blockindex"`
	BlockTime int                `json:"blocktime"`
	TxId string                  `json:"txid"`
	Time int                     `json:"time"`
	TimeReceived int             `json:"timereceived"`
	Bip125Replaceable string     `json:"bip125-replaceable"`
	Details []TransactionDetails `json:"details"`
	Hex           string         `json:"hex"`
}

type TransactionDetails struct {
	Address string       `json:"address"`
	Category string       `json:"category"`
	Amount float64       `json:"amount"`
	Label string       `json:"label"`
	Vout int       `json:"vout"`
	Fee float64       `json:"fee"`
	Abandoned bool       `json:"abandoned"`
}

type RawTransaction struct {
	InActiveChain bool      `json:"in_active_chain"`
	Hex string      `json:"hex"`
	TxId string      `json:"txid"`
	Hash string      `json:"hash"`
	Size int64      `json:"size"`
	Vsize int64      `json:"vsize"`
	Weight int `json:"weight"`
	LockTime int      `json:"locktime"`
	Version int      `json:"version"`
	BlockHash string       `json:"blockhash"`
	BlockTime int       `json:"blocktime"`
	Confirmations int       `json:"confirmations"`
	Time int       `json:"time"`
	Vout []struct{
		Amount float64 `json:"value"`
		N int `json:"n"`
		ScriptPubKey struct{
			Asm string `json:"asm"`
			Hex string`json:"hex"`
			ReqSigs int`json:"regSigs"`
			Type string`json:"type"`
			Addresses []string`json:"addresses"`
		} `json:"scriptPubKey"`
	}
	Vin []struct{
		TxId string `json:"txid"`
		Vout int `json:"vout"`
		ScriptSig struct{
			Asm string	`json:"asm"`
			Hex string	`json:"hex"`
		}	`json:"scriptSig"`
		Sequence int	`json:"sequence"`
		TxInWitness []string `json:"txinwitness"`
	}
}

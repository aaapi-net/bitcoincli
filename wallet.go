package bitcoincli

type Wallet struct {
	Name string       `json:"name"`
}

type WalletWithWarning struct {
	Wallet
	Warning string      `json:"warning"`
}

type ListWalletDir struct {
	Wallets []Wallet `json:"wallets"`
}

package model

type BalanceBody struct {
	Address            string `json:"address"`
	TotalReceived      int64  `json:"total_received"`
	TotalSent          int64  `json:"total_sent"`
	Balance            int64  `json:"balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	FinalBalance       int64  `json:"final_balance"`
	NTx                int    `json:"n_tx"`
	UnconfirmedNTx     int    `json:"unconfirmed_n_tx"`
	FinalNTx           int    `json:"final_n_tx"`
}

type GetInfo struct {
	Difficulty int `json:"Difficulty"`
	Hashrate   int `json:"Hashrate"`
	CoinSupply int `json:"CoinSupply"`
}

type GetBalance struct {
	RAIDS int `json:"RAIDS"`
	BTC   int `json:"BTC"`
}

type GetNewAddress struct {
	DeviceId  string `json:"deviceId"`
	Timestamp int    `json:"timestamp"`
	Hmac      string `json:"hmac"`
}

type GetNewAddressResponse struct {
	Address    string `json:"Address"`
	PrivateKey string `json:"PrivateKey"`
}

type GetAuthToken struct {
	XAuthToken string `json:"XAuthToken"`
}

type Block struct {
	Height int64
	Timestamp int
	RelayedBy string
	Difficulty int
	HashResult string
	BlockReward float64
}
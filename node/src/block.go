package src

type Block struct {
	ID           string         `json:"ID"`
	Balances     map[string]int `json:"Balances"`
	Transactions []string       `json:"Transactions"`
	Target       int            `json:"target"`
	Nonce        int            `json:"nonce"`
	Previous     *Block         `json:"Previous" copy:"shallow"`
	Next         *Block         `json:"next" copy:"shallow"`
}
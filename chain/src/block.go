package src

type Block struct {
	ID string
	Balances map[string]int
	Transactions []string
	Target int
	Nonce int
	Previous *Block
	Next *Block
}
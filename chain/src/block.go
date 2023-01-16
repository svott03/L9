package src

type Block struct {
	ID           string
	Balances     map[string]int
	Transactions []string
	Target       int
	Nonce        int
	Previous     *Block
	Next         *Block
}

type sampleBlock struct {
	ID string
}
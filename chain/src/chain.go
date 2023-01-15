package src

import (
	"net/http"
	"encoding/binary"
)

type Chain struct {
	NodeCount uint32;
	CoinCount uint32;
	Root *Block;
}

func (c *Chain) Run() {
	c.NodeCount = 0;
	c.Root = &Block{
		ID: "Genesis",
	}
	// Listens to Concurrent requests
	http.HandleFunc("/join", join(c))
	http.HandleFunc("/verify", verify(c))

	http.ListenAndServe(":8080", nil)

}

// Incoming Nodes
func join(c *Chain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs := make([]byte, 4)
    binary.LittleEndian.PutUint32(bs, c.NodeCount)
		w.Write(bs)
		c.NodeCount++
	}
}

// Incoming blocks
func verify(c *Chain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// verify hash

		// will only add block that contains the current block as previous block
	}
}



// maintain a linked list of blocks

// pool of transactions

// TODO encode transactions in blocks
// TODO validate transactions

// TODO while broadcasting, cannot accept new incoming requests
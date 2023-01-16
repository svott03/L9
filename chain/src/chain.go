package src

import (
	"net/http"
	"encoding/binary"
	"encoding/json"
	"bytes"
	"strconv"
	"fmt"
	"cmp"
)

type Chain struct {
	NodeCount uint32;
	CoinCount uint32;
	BlockCount uint32;
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
		w.Header().Set("Content-Type", "application/json")
		var data struct {
			Status string
			CoinReward int
		}
		// verify hash

		// will only add block that contains the current block as previous block
		var receivedBlock Block;
		err := json.NewDecoder(r.Body).Decode(&receivedBlock)
		if err != nil {
			fmt.Println(err)
		}
		var jData []byte
		// May panic if objects are not both blocks, its ok for this project though
		if cmp.Equal(receivedBlock.Previous,c.Root) {
			// TODO mux
			c.broadCast();
			data.Status = "Block Accepted!"
			data.CoinReward = 1
			jData, _ = json.Marshal(data)
		} else {
			data.Status = "You are not up to date with the chain"
			data.CoinReward = 0
			jData, _ = json.Marshal(data)
		}
		w.Write(jData)
	}
}

// Broadcast new block to all nodes connected
func (c* Chain) broadCast() {
	data, _ := json.Marshal(c.Root)
	b := bytes.NewBuffer(data)
	for i := uint32(0); i < c.NodeCount; i++ {
		_, err := http.Post("http://localhost:" + strconv.Itoa(int(c.NodeCount)) + "/newBlock", "application/json", b)
		if err != nil {
			fmt.Println(err)
		}
	}
}


// maintain a linked list of blocks

// pool of transactions

// TODO encode transactions in blocks
// TODO validate transactions

// TODO while broadcasting, cannot accept new incoming requests
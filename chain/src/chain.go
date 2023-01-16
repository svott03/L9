package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"strconv"
)

type Chain struct {
	NodeCount  uint32
	CoinCount  uint32
	BlockCount uint32
	Root       *Block
}

func (c *Chain) Run() {
	c.NodeCount = 0
	c.Root = &Block{
		ID:       "Genesis",
		Balances: make(map[string]int),
	}
	// Listens to Concurrent requests
	http.HandleFunc("/join", join(c))
	http.HandleFunc("/verify", verify(c))

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}

// Incoming Nodes
func join(c *Chain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Port      uint32
			RootBlock *Block
		}
		w.Header().Set("Content-Type", "application/json")
		data.Port = c.NodeCount
		data.RootBlock = c.Root
		encoder := json.NewEncoder(w)
		encoder.Encode(data)
		c.NodeCount++
		color.Green("Nodes connected: " + strconv.Itoa(int(c.NodeCount)))
	}
}

// Incoming blocks
func verify(c *Chain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var data struct {
			Status     string
			CoinReward int
		}
		var jData []byte
		// verify hash

		// will only add block that contains the current block as previous block
		var receivedBlock Block
		err := json.NewDecoder(r.Body).Decode(&receivedBlock)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", receivedBlock)
		// May panic if objects are not both blocks
		// fmt.Println(receivedBlock.Previous)
		// fmt.Println(c.Root)
		if cmp.Equal(receivedBlock.Previous, c.Root) {
			// TODO mux
			c.broadCast()
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
func (c *Chain) broadCast() {
	fmt.Println("Broadcasting...")
	data, _ := json.Marshal(c.Root)
	b := bytes.NewBuffer(data)
	for i := uint32(0); i < c.NodeCount; i++ {
		_, err := http.Post("http://localhost:"+strconv.Itoa(int(c.NodeCount))+"/newBlock", "application/json", b)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// TODO validate transactions

// TODO while broadcasting, cannot accept new incoming requests

// TODO hashing

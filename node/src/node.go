package src

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"bytes"
)

type Node struct {
	port uint32
	CurBlock *Block
	miningAttempts int
	coins int
}

// TODO put mutex on Node's block
func (n *Node) Run() {
	// Connect to blockchain
	res, err := http.Get("http://localhost:8080/join")
	if err != nil {
		fmt.Println(err)
	}
	var r struct {
		Port uint32
		RootBlock Block
	}
	_ = json.NewDecoder(res.Body).Decode(&r)
	n.port = r.Port
	// GC transfers ownership, so the deallocation of r does not affect n.CurBlock
	n.CurBlock = &r.RootBlock
	fmt.Println("Connected to Blockchain! Port is " + strconv.Itoa(int(n.port)))

	// need routine for I/O
	go Input(n)

	// listen on port
	http.HandleFunc("/newBlock", newBlock(n))
	http.ListenAndServe(":" + strconv.Itoa(int(n.port)), nil)
}

func Input(n *Node) {
	for {
		// accept input
		fmt.Println("----------------------")
		fmt.Println("Select an option below")
		fmt.Println("1. Make Transaction")
		fmt.Println("2. Mine Coin")
		// handle input
		var option int
		argc, err := fmt.Scanln(&option)
		if argc != 1 || err != nil || (option != 1 && option != 2) {
			fmt.Println("Please output either 1 or 2.")
			continue
		}
		if option == 1 {
			n.transact()
		} else {
			n.mine()
		}
	}
}

// Updates Node's block to the Chain's latest block
func newBlock(n *Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// will only add block that contains the current block as previous block
		err := json.NewDecoder(r.Body).Decode(n.CurBlock)
		n.miningAttempts = 0
		if err != nil {
			fmt.Println(err)
		}
	}
}

// TODO lock node
// Incoming Transactions
func (n *Node) transact() {
	var operation string
	var name1 string
	var name2 string
	var balance int
	// INSERT balance name1
	// TRADE balance name1 name2
	fmt.Println("Please enter a transaction following the format")
	fmt.Println("INSERT <balance> <name1> will add a user to the ledger")
	fmt.Println("TRADE <balance> <name1> <name2> will transfer <balance> funds from <name1>'s account to <name2>'s account")
	fmt.Scanln(&operation, &balance, &name1, &name2)
	// execute trade on current block
	if operation == "INSERT" {
		n.CurBlock.Balances[name1] += balance
		n.CurBlock.Transactions = append(n.CurBlock.Transactions, "Inserted " + strconv.Itoa(balance) + " funds into account " + name1)
	} else {
		if n.CurBlock.Balances[name1] < balance {
			fmt.Println("Could not execute transaction. Insufficient funds.")
		} else {
			n.CurBlock.Balances[name1] -= balance
			n.CurBlock.Balances[name2] += balance
			n.CurBlock.Transactions = append(n.CurBlock.Transactions, "Transfered " + strconv.Itoa(balance) + " funds from " + name1 + " to " + name2)
		}
	}
	fmt.Println("Transaction Complete.")
}

func (n *Node) mine() {
	fmt.Println("Mining attempt " + strconv.Itoa(n.miningAttempts))
	data, _ := json.Marshal(n.CurBlock)
	b := bytes.NewBuffer(data)
	res, err := http.Post("http://localhost:8080/verify", "application/json", b)
	if err != nil {
		fmt.Println(err)
	}
	var response struct {
		Status string
		CoinReward int
	}
	_ = json.NewDecoder(res.Body).Decode(&response)
	if response.CoinReward > 0 {
		fmt.Println("Success. " + response.Status)
		n.coins += response.CoinReward
	} else {
		fmt.Println("Failure. " + response.Status)
	}
	n.miningAttempts++
}

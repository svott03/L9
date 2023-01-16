package src

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
)

type Node struct {
	ID    int
	Port int
	CurBlock *Block
}

// TODO put mutex on Node's block

func (n *Node) Run() {
	// Connect to blockchain
	res, err := http.Get("http://localhost:8080/join")
	if err != nil {
		fmt.Println(err)
	}
	var r struct {
		Port int
	}
	decoder := json.NewDecoder(res.Body)
	_ = decoder.Decode(&r)
	n.Port = r.Port
	fmt.Println("Connected to Blockchain! Port is " + strconv.Itoa(n.Port))

	// listen on port
	http.HandleFunc("/newBlock", newBlock(n))
	http.ListenAndServe(":" + strconv.Itoa(n.Port), nil)

	fmt.Println("Connected to L9 blockchain with id: " + strconv.Itoa(n.ID))
	// listen to user input

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
	} else {
		if n.CurBlock.Balances[name1] < balance {
			fmt.Println("Could not execute transaction. Insufficient funds.")
		} else {
			n.CurBlock.Balances[name1] -= balance
			n.CurBlock.Balances[name2] += balance
		}
	}
	fmt.Println("Transaction Complete.")
}

func (n *Node) mine() {
	// output mining attempts

	// make post request to chain's /verify
}

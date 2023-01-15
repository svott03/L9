package src

import (
	"fmt"
	"net/http"
	"strconv"
)

type Node struct {
	ID    int
	block *Block
}

// TODO put mutex on Node's block

func (n *Node) Run() {
	var port int
	// Connect to blockchain

	// listen on port
	http.HandleFunc("/newBlock", join(c))
	http.ListenAndServe(":" + strconv.Itoa(port), nil)

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

	} else {

	}
}

func (n *Node) mine() {
	// output mining attempts
}

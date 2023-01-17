package src

import (
	"bytes"
	// "reflect"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/mitchellh/copystructure"
)

type Node struct {
	port           uint32
	CurBlock       *Block
	ChainBlock     *Block `json:"chain_block"`
	miningAttempts int
	coins          int
}

func print_list(root *Block) {
	if root == nil {return}
	fmt.Printf("%+v\n", *root)
	print_list(root.Next)
}

// TODO put mutex on Node's block
func (n *Node) Run() {
	// Connect to blockchain
	res, err := http.Get("http://localhost:8080/join")
	if err != nil {
		fmt.Println(err)
	}
	var r struct {
		Port      uint32
		RootBlock *Block
	}
	_ = json.NewDecoder(res.Body).Decode(&r)
	n.port = r.Port
	// GC transfers ownership
	n.ChainBlock = r.RootBlock
	t, _ := copystructure.Copy(n.ChainBlock)
	n.CurBlock = t.(*Block)
	// n.CurBlock = reflect.ValueOf(t).Interface().(*Block)
	n.CurBlock.Previous = n.ChainBlock
	// n.ChainBlock.Next = n.CurBlock
	print_list(n.ChainBlock)
	fmt.Println("Connected to Blockchain! Port is " + strconv.Itoa(int(n.port)))
	// need routine for I/O
	go Input(n)
	// listen on port
	http.HandleFunc("/newBlock", newBlock(n))
	http.ListenAndServe(":"+strconv.Itoa(int(n.port)), nil)
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
		fmt.Println("Accepting New Block...")
		// will only add block that contains the current block as previous block
		err := json.NewDecoder(r.Body).Decode(n.ChainBlock)
		n.CurBlock = n.ChainBlock
		n.CurBlock.Previous = n.ChainBlock
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
		n.CurBlock.Transactions = append(n.CurBlock.Transactions, "Inserted "+strconv.Itoa(balance)+" funds into account "+name1)
	} else {
		if n.CurBlock.Balances[name1] < balance {
			fmt.Println("Could not execute transaction. Insufficient funds.")
		} else {
			n.CurBlock.Balances[name1] -= balance
			n.CurBlock.Balances[name2] += balance
			n.CurBlock.Transactions = append(n.CurBlock.Transactions, "Transfered "+strconv.Itoa(balance)+" funds from "+name1+" to "+name2)
		}
	}
	fmt.Println("Transaction Complete.")
}




func (n *Node) mine() {
	fmt.Println("Mining attempt " + strconv.Itoa(n.miningAttempts))
	data, err3 := json.Marshal(*(n.CurBlock))
	if err3 != nil {
		fmt.Println(err3)
	}
	b := bytes.NewBuffer(data)
	fmt.Println("----------List-------------")
	print_list(n.CurBlock)
	fmt.Println("------------Data----------")
	fmt.Println(data)
	fmt.Println(string(data))
	fmt.Println("---------------Mine---------------")
	fmt.Println(n.CurBlock)
	fmt.Println(n.ChainBlock)
	fmt.Println(n.ChainBlock.Next)
	fmt.Println(n.CurBlock.Previous)

	var receivedBlock *Block
	err := json.NewDecoder(b).Decode(&receivedBlock)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", receivedBlock)

	res, err := http.Post("http://localhost:8080/verify", "application/json", b)
	if err != nil {
		fmt.Println(err)
	}
	var response struct {
		Status     string
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

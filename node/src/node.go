package src

import (
	"net/http"
)

type Node struct {
	ID int
}



// Incoming Transactions
func transact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// transaction logic
	}
}
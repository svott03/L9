# L9 BlockChain 🪙
A naive implementation of a blockchain.

In our implementation, each block holds balances and a set of transactions made.

Nodes will try to import only feasible transactions, that is nodes will check previous block data prevent double spending.

When a new block is added, we must lock all other incoming block requests and update the chain. We then distribute the latest state of the chain to all nodes and the process starts all over again.

## Discontinued
We discontinue this project as cannot nicely serialize recursive data structures. 

## TODO
Authentication
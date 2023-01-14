package src

type Chain struct {
	NodeCount int;
	Root *Block;
}

func (c *Chain) Run() {
	c.NodeCount = 0;
	c.Root = &Block{
		ID: "Genesis",
	}
	// listen to endpoints
}

func (c *Chain) nodeAssign {
	// assign node id
}


// maintain a linked list of blocks

// pool of transactions

// TODO encode transactions in blocks
// TODO validate transactions
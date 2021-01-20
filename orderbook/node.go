package orderbook

import rbt "github.com/emirpasic/gods/trees/redblacktree"

type node struct {
	node *rbt.Node
}

func (n *node) amount() float64 {
	return n.node.Value.(float64)
}

func (n *node) setAmount(amount float64) {
	n.node.Value = amount
}

func (n *node) price() float64 {
	return n.node.Key.(float64)
}

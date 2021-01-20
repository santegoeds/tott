package orderbook

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
)

type tree struct {
	tree *rbt.Tree
}

func newSellTree() *tree {
	return &tree{rbt.NewWith(utils.Float64Comparator)}
}

func newBuyTree() *tree {
	return &tree{rbt.NewWith(func(a, b interface{}) int {
		return utils.Float64Comparator(b, a)
	})}
}

func (t *tree) Get(price float64) (float64, bool) {
	v, ok := t.tree.Get(price)
	if !ok {
		return 0.0, false
	}
	return v.(float64), true
}

func (t *tree) Put(price float64, amount float64) {
	t.tree.Put(price, amount)
}

func (t *tree) Remove(price float64) {
	t.tree.Remove(price)
}

func (t *tree) floor(price float64) (*node, bool) {
	n, ok := t.tree.Floor(price)
	if !ok {
		return nil, false
	}
	return &node{n}, true
}

func (t *tree) ceiling(price float64) (*node, bool) {
	n, ok := t.tree.Ceiling(price)
	if !ok {
		return nil, false
	}
	return &node{n}, true
}

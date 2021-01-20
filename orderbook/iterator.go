package orderbook

import rbt "github.com/emirpasic/gods/trees/redblacktree"

type Iterator struct {
	it rbt.Iterator
}

func (t *tree) Iterator() Iterator {
	return Iterator{t.tree.Iterator()}
}

func (it *Iterator) Next() bool {
	return it.it.Next()
}

func (it *Iterator) Prev() bool {
	return it.it.Prev()
}

func (it *Iterator) First() bool {
	return it.it.First()
}

func (it *Iterator) Last() bool {
	return it.it.Last()
}

func (it *Iterator) Price() float64 {
	return it.it.Key().(float64)
}

func (it *Iterator) Amount() float64 {
	return it.it.Value().(float64)
}

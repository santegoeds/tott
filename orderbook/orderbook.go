package orderbook

import (
	"github.com/kr/pretty"
	"log"
	"math"
)

type Side string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

type PriceLevel struct {
	Price  float64
	Amount float64
}

type Front struct {
	Buy  PriceLevel
	Sell PriceLevel
}

func (f Front) Spread() float64 {
	return f.Sell.Price - f.Buy.Price
}

type Order interface {
	ID() string
	Price() float64
	Amount() float64
	Side() Side
}

type PricerFunc func(o Order) float64

type Orderbook struct {
	orders  map[string]Order
	buys    *tree
	sells   *tree
	priceFn PricerFunc
}

func New() *Orderbook {
	ob := &Orderbook{
		orders:  make(map[string]Order),
		buys:    newBuyTree(),
		sells:   newSellTree(),
		priceFn: nil,
	}
	return ob
}

func NewWithPricer(priceFn PricerFunc) *Orderbook {
	ob := New()
	ob.priceFn = priceFn
	return ob
}

func NewWithPrecision(precision float64) *Orderbook {
	if precision <= 0.0 {
		return NewWithPricer(func(o Order) float64 {
			return o.Price()
		})
	}

	// Rounds prices to the specified precision.  Buy orders are rounded up
	// and Sell orders are rounded down so that slippage of orders based on
	// the orderbook have a positive bias.
	return NewWithPricer(func(o Order) float64 {
		price := o.Price()
		rem := math.Remainder(price, precision)
		if rem == 0.0 {
			return price
		}
		if o.Side() == Buy {
			return price - rem + precision
		}
		return price - rem
	})
}

func (b *Orderbook) Add(o Order) {
	b.orders[o.ID()] = o

	var tree *tree
	if o.Side() == Buy {
		tree = b.buys
	} else {
		tree = b.sells
	}

	var price float64
	if b.priceFn == nil {
		price = o.Price()
	} else {
		price = b.priceFn(o)
	}

	node, ok := tree.floor(price)
	if ok && node.price() == price {
		// A node for this price already exists.
		node.setAmount(node.amount() + o.Amount())
		return
	}
	tree.Put(price, o.Amount())
}

func (b *Orderbook) Remove(id string) {
	// Remove the order from the index.
	o, ok := b.orders[id]
	if !ok {
		return
	}
	delete(b.orders, id)

	// Remove the order from the tree.
	var tree *tree
	if o.Side() == Buy {
		tree = b.buys
	} else {
		tree = b.sells
	}

	var price float64
	if b.priceFn == nil {
		price = o.Price()
	} else {
		price = b.priceFn(o)
	}

	node, ok := tree.floor(price)
	if !ok {
		log.Printf("Order %s not found in tree", pretty.Sprint(o))
		return
	}
	if node.price() < price {
		log.Printf("Block for Order %s not found in tree", pretty.Sprint(o))
		return
	}

	if node.amount() <= o.Amount() {
		// Remove the block if this was the last order.
		tree.Remove(price)
	} else {
		// Decrease the amount in the order's block.
		node.setAmount(node.amount() - o.Amount())
	}
}

func (b *Orderbook) Iterator(side Side) Iterator {
	var tree *tree
	if side == Buy {
		tree = b.buys
	} else {
		tree = b.sells
	}
	return tree.Iterator()
}

func (b *Orderbook) Front(amount float64) Front {
	return Front{
		Buy:  b.BestSide(Buy, amount),
		Sell: b.BestSide(Sell, amount),
	}
}

func (b *Orderbook) BestSide(side Side, amount float64) PriceLevel {
	var tree *tree
	if side == Buy {
		tree = b.buys
	} else {
		tree = b.sells
	}

	it := tree.Iterator()
	if ok := it.First(); !ok {
		return PriceLevel{}
	}
	// Front of the book
	if amount == 0.0 {
		return PriceLevel{Price: it.Price(), Amount: it.Amount()}
	}

	totalAmount := it.Amount()
	totalPrice := it.Price() * it.Amount()
	amount -= math.Min(amount, it.Amount())
	for amount > 0 && it.Next() {
		neededAmount := math.Min(amount, it.Amount())
		totalAmount += neededAmount
		totalPrice += it.Price() * neededAmount
		amount -= neededAmount
	}
	return PriceLevel{Price: totalPrice / totalAmount, Amount: totalAmount}
}

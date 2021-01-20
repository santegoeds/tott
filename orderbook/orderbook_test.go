package orderbook_test

import (
	"math"
	"testing"

	"github.com/santegoeds/tott/orderbook"
	"github.com/stretchr/testify/require"
)

type Order struct {
	id     string
	price  float64
	amount float64
	side   orderbook.Side
}

func (o Order) ID() string {
	return o.id
}

func (o Order) Price() float64 {
	return o.price
}

func (o Order) Amount() float64 {
	return o.amount
}

func (o Order) Side() orderbook.Side {
	return o.side
}

func TestOrderbook(t *testing.T) {
	ob := orderbook.New()
	buys := []Order{
		{id: "1", price: 0.1, amount: 100, side: orderbook.Buy},
		{id: "2", price: 0.2, amount: 100, side: orderbook.Buy},
	}
	sells := []Order{
		{id: "3", price: 0.4, amount: 150, side: orderbook.Sell},
		{id: "4", price: 0.3, amount: 50, side: orderbook.Sell},
		{id: "5", price: 0.3, amount: 50, side: orderbook.Sell},
	}
	for _, o := range append(buys, sells...) {
		ob.Add(o)
	}
	front := ob.Front(0)
	require.Equal(t, orderbook.Front{
		Buy:  orderbook.PriceLevel{Price: 0.2, Amount: 100.0},
		Sell: orderbook.PriceLevel{Price: 0.3, Amount: 100.0},
	}, front)
	require.InDelta(t, 0.1, front.Spread(), 1e-16)

	front = ob.Front(200)
	require.Equal(t, orderbook.Front{
		Buy:  orderbook.PriceLevel{Price: 0.15, Amount: 200.0},
		Sell: orderbook.PriceLevel{Price: 0.35, Amount: 200.0},
	}, front)
	require.InDelta(t, 0.2, front.Spread(), 1e-16)

	prevPrice := math.Inf(-1)
	it := ob.Iterator(orderbook.Buy)
	for it.Next() {
		require.Greater(t, it.Amount(), 0.0)
		require.Less(t, prevPrice, it.Price())
	}
	prevPrice = math.Inf(1)
	it = ob.Iterator(orderbook.Sell)
	for it.Next() {
		require.Greater(t, it.Amount(), 0.0)
		require.Greater(t, prevPrice, it.Price())
	}

	for _, o := range append(buys, sells...) {
		ob.Remove(o.ID())
	}
	front = ob.Front(100)
	require.Equal(t, orderbook.Front{
		Buy:  orderbook.PriceLevel{Price: 0.0, Amount: 0.0},
		Sell: orderbook.PriceLevel{Price: 0.0, Amount: 0.0},
	}, front)
}

package db

import (
	"github.com/google/uuid"
	"log"
	"orderbook/db/models"
)

type ActionType string

const (
	AtBuy           = "BUY"
	AtSell          = "SELL"
	AtPartialFilled = "PARTIAL_FILLED"
	AtFilled        = "FILLED"
)

const MaxPrice = 10000000

type PricePoint struct {
	OrderHead *models.Order
	OrderTail *models.Order
}

func (pp *PricePoint) Insert(o *models.Order) {
	if pp.OrderHead == nil {
		pp.OrderHead = o
		pp.OrderTail = o
	} else {
		pp.OrderTail.Next = o
		pp.OrderTail = o
	}
}

type OrderBook struct {
	MinAsk      uint32
	MaxBid      uint32
	OrderIndex  map[uuid.UUID]*models.Order
	PricePoints [MaxPrice]*PricePoint
}

func NewOrderBookMap() *[]*OrderBook {
	var obMap []*OrderBook
	symbols := []models.Symbol{models.BTCUSDT, models.ETHUSDT, models.BTCIRT, models.ETHIRT}
	for _, s := range symbols {
		obMap[s] = newOrderBook()
	}

	return &obMap
}

func newOrderBook() *OrderBook {
	ob := new(OrderBook)
	ob.MaxBid = 0
	ob.MinAsk = MaxPrice
	for i := range ob.PricePoints {
		ob.PricePoints[i] = new(PricePoint)
	}
	ob.OrderIndex = make(map[uuid.UUID]*models.Order)
	return ob
}

func (ob *OrderBook) AddOrder(o *models.Order) {
	if o.Side == models.Buy {
		log.Printf("actionType: %s, order: %s", AtBuy, o.String())
		ob.FillBuy(o)
	} else {
		log.Printf("actionType: %s, order: %s", AtSell, o.String())
		ob.FillSell(o)
	}

	if o.Amount > 0 {
		ob.openOrder(o)
	}
}

func (ob *OrderBook) openOrder(o *models.Order) {
	pp := ob.PricePoints[o.Price]
	pp.Insert(o)
	o.Status = models.OsOpen
	if o.Side == models.Buy && o.Price > ob.MaxBid {
		ob.MaxBid = o.Price
	} else if o.Side != models.Buy && o.Price < ob.MinAsk {
		ob.MinAsk = o.Price
	}
	ob.OrderIndex[o.ID] = o
}

func (ob *OrderBook) FillBuy(o *models.Order) {
	for ob.MinAsk <= o.Price && o.Amount > 0 {
		pp := ob.PricePoints[ob.MinAsk]
		ppOrderHead := pp.OrderHead
		for ppOrderHead != nil {
			ob.fill(o, ppOrderHead)
			ppOrderHead = ppOrderHead.Next
			pp.OrderHead = ppOrderHead
		}
		ob.MinAsk++
	}
}

func (ob *OrderBook) FillSell(o *models.Order) {
	for ob.MaxBid >= o.Price && o.Amount > 0 {
		pp := ob.PricePoints[ob.MaxBid]
		ppOrderHead := pp.OrderHead
		for ppOrderHead != nil {
			ob.fill(o, ppOrderHead)
			ppOrderHead = ppOrderHead.Next
			pp.OrderHead = ppOrderHead
		}
		ob.MaxBid--
	}
}

func (ob *OrderBook) fill(o, ppOrderHead *models.Order) {
	if ppOrderHead.Amount >= o.Amount {
		log.Printf("actionType: %s, order: %s, fromOrder: %s", AtFilled, o.String(), ppOrderHead.String())
		ppOrderHead.Amount -= o.Amount
		o.Amount = 0
		o.Status = models.OsFilled
		return
	} else {
		if ppOrderHead.Amount > 0 {
			log.Printf("actionType: %s, order: %s, fromOrder: %s", AtPartialFilled, o.String(), ppOrderHead.String())
			o.Amount -= ppOrderHead.Amount
			o.Status = models.OsPartial
			ppOrderHead.Amount = 0
		}
	}
}

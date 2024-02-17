package core

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	apimodels "orderbook/api/models"
	"orderbook/core/models"
)

type ActionType string

const (
	AtBuy           = "BUY"
	AtSell          = "SELL"
	AtPartialFilled = "PARTIAL_FILLED"
	AtFilled        = "FILLED"
)

const MaxPrice = 100000000

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

type Orderbook struct {
	MinAsk      uint32
	MaxBid      uint32
	OrderIndex  map[uuid.UUID]*models.Order
	PricePoints [MaxPrice]*PricePoint
}

var OrderbookMap *map[models.Symbol]*Orderbook

func NewOrderbookMap() {
	obMap := make(map[models.Symbol]*Orderbook)
	symbols := []models.Symbol{models.BTCUSDT, models.ETHUSDT, models.BTCIRT, models.ETHIRT}
	for _, s := range symbols {
		obMap[s] = newOrderbook()
	}

	OrderbookMap = &obMap
}

func newOrderbook() *Orderbook {
	ob := new(Orderbook)
	ob.MaxBid = 0
	ob.MinAsk = MaxPrice
	for i := range ob.PricePoints {
		ob.PricePoints[i] = new(PricePoint)
	}
	ob.OrderIndex = make(map[uuid.UUID]*models.Order)
	return ob
}

func (ob *Orderbook) AddOrder(o *models.Order) {
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

func (ob *Orderbook) openOrder(o *models.Order) {
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

func (ob *Orderbook) FillBuy(o *models.Order) {
	for ob.MinAsk <= o.Price && o.Amount > 0 {
		pp := ob.PricePoints[ob.MinAsk]
		ppOrderHead := pp.OrderHead
		for ppOrderHead != nil {
			ob.fill(o, ppOrderHead)
			if o.Amount == 0 {
				return
			}
			ppOrderHead = ppOrderHead.Next
			pp.OrderHead = ppOrderHead
		}
		ob.MinAsk++
	}
}

func (ob *Orderbook) FillSell(o *models.Order) {
	for ob.MaxBid >= o.Price && o.Amount > 0 {
		pp := ob.PricePoints[ob.MaxBid]
		ppOrderHead := pp.OrderHead
		for ppOrderHead != nil {
			ob.fill(o, ppOrderHead)
			if o.Amount == 0 {
				return
			}
			ppOrderHead = ppOrderHead.Next
			pp.OrderHead = ppOrderHead
		}
		ob.MaxBid--
	}
}

func (ob *Orderbook) fill(o, ppOrderHead *models.Order) {
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

func (ob *Orderbook) GetBidsTable(maxN int) []apimodels.OrderRecord {
	var bids []apimodels.OrderRecord

	for mb := ob.MaxBid; len(bids) < maxN && mb > 0; mb-- {
		pp := ob.PricePoints[mb]
		o := pp.OrderHead
		total := 0.
		for o != nil {
			total += o.Amount
			o = o.Next
		}
		if total > 0 {
			bids = append(bids, apimodels.OrderRecord{fmt.Sprintf("%d", mb), fmt.Sprintf("%f", total)})
		}
	}
	return bids
}

func (ob *Orderbook) GetAsksTable(maxN int) []apimodels.OrderRecord {
	var asks []apimodels.OrderRecord

	for ma := ob.MinAsk; len(asks) < maxN && ma < MaxPrice; ma++ {
		pp := ob.PricePoints[ma]
		o := pp.OrderHead
		total := 0.
		for o != nil {
			total += o.Amount
			o = o.Next
		}
		if total > 0 {
			asks = append(asks, apimodels.OrderRecord{fmt.Sprintf("%d", ma), fmt.Sprintf("%f", total)})
		}
	}
	return asks
}

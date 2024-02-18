package core

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	apimodels "orderbook/api/models"
	"orderbook/core/models"
	"orderbook/db/mongo"
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
	OrderIndex  map[primitive.ObjectID]*models.Order
	PricePoints [MaxPrice]*PricePoint
}

var OrderbookMap *map[models.Symbol]*Orderbook

func Init() {
	obMap := make(map[models.Symbol]*Orderbook)
	symbols := []models.Symbol{models.BTCUSDT, models.ETHUSDT, models.BTCIRT, models.ETHIRT}
	for _, s := range symbols {
		obMap[s] = newOrderbook()
	}
	OrderbookMap = &obMap

	addMongoOrders()
}

func addMongoOrders() {
	all, _ := mongo.OrderRepo.GetByFilter(bson.D{{"remain_amount", bson.D{{"$gt", 0}}}})

	for _, o := range all {
		(*OrderbookMap)[o.Symbol].openOrder(o)
	}
	log.Println(len(all), "orders added to orderbooks!")
}

func newOrderbook() *Orderbook {
	ob := new(Orderbook)
	ob.MaxBid = 0
	ob.MinAsk = MaxPrice
	for i := range ob.PricePoints {
		ob.PricePoints[i] = new(PricePoint)
	}
	ob.OrderIndex = make(map[primitive.ObjectID]*models.Order)
	return ob
}

func (ob *Orderbook) AddOrder(o *models.Order) {
	_ = mongo.OrderRepo.Create(o)

	if o.Side == models.Buy {
		log.Printf("start %s process for %s", AtBuy, o.String())
		ob.FillBuy(o)
	} else {
		log.Printf("start %s process for %s", AtSell, o.String())
		ob.FillSell(o)
	}

	if o.RemainAmount > 0 {
		ob.openOrder(o)
	}

	_ = mongo.OrderRepo.Update(o)
}

func (ob *Orderbook) openOrder(o *models.Order) {
	pp := ob.PricePoints[o.Price]
	pp.Insert(o)
	if o.Side == models.Buy && o.Price > ob.MaxBid {
		ob.MaxBid = o.Price
	} else if o.Side != models.Buy && o.Price < ob.MinAsk {
		ob.MinAsk = o.Price
	}
	ob.OrderIndex[o.ID] = o
}

func (ob *Orderbook) FillBuy(o *models.Order) {
	for ob.MinAsk <= o.Price && o.RemainAmount > 0 {
		pp := ob.PricePoints[ob.MinAsk]
		ppOrderHead := pp.OrderHead
		for ppOrderHead != nil {
			ob.fill(o, ppOrderHead)
			if o.RemainAmount == 0 {
				return
			}
			ppOrderHead = ppOrderHead.Next
			pp.OrderHead = ppOrderHead
		}
		ob.MinAsk++
	}
}

func (ob *Orderbook) FillSell(o *models.Order) {
	for ob.MaxBid >= o.Price && o.RemainAmount > 0 {
		pp := ob.PricePoints[ob.MaxBid]
		ppOrderHead := pp.OrderHead
		for ppOrderHead != nil {
			ob.fill(o, ppOrderHead)
			if o.RemainAmount == 0 {
				return
			}
			ppOrderHead = ppOrderHead.Next
			pp.OrderHead = ppOrderHead
		}
		ob.MaxBid--
	}
}

func (ob *Orderbook) fill(o *models.Order, from *models.Order) {
	if from.RemainAmount >= o.RemainAmount {
		log.Printf("%s %s from %s", o.String(), AtFilled, from.String())
		from.RemainAmount -= o.RemainAmount
		o.RemainAmount = 0
	} else {
		if from.RemainAmount > 0 {
			log.Printf("%s %s from %s", o.String(), AtPartialFilled, from.String())
			o.RemainAmount -= from.RemainAmount
			from.RemainAmount = 0
		}
	}
	_ = mongo.OrderRepo.Update(o)
	_ = mongo.OrderRepo.Update(from)
}

func (ob *Orderbook) GetBidsTable(maxN int) []apimodels.OrderRecord {
	var bids []apimodels.OrderRecord

	for mb := ob.MaxBid; len(bids) < maxN && mb > 0; mb-- {
		pp := ob.PricePoints[mb]
		o := pp.OrderHead
		total := 0.
		for o != nil {
			total += o.RemainAmount
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
			total += o.RemainAmount
			o = o.Next
		}
		if total > 0 {
			asks = append(asks, apimodels.OrderRecord{fmt.Sprintf("%d", ma), fmt.Sprintf("%f", total)})
		}
	}
	return asks
}

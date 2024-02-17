package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Symbol string

const (
	BTCUSDT Symbol = "BTCUSDT"
	ETHUSDT Symbol = "ETHUSDT"
	BTCIRT  Symbol = "BTCIRT"
	ETHIRT  Symbol = "ETHIRT"
)

type Side string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

type OrderStatus int

const (
	OsNew OrderStatus = iota
	OsOpen
	OsPartial
	OsFilled
)

type Order struct {
	ID     uuid.UUID   `json:"id"`
	Side   Side        `json:"side"`
	Symbol Symbol      `json:"symbol"`
	Amount float64     `json:"amount"`
	Price  uint32      `json:"price"`
	Status OrderStatus `json:"status"`
	Next   *Order      `json:"next"`
}

func (o *Order) String() string {
	return fmt.Sprintf("Order{id:%v,isBuy:%v,price:%v,amount:%v}", o.ID, o.Side, o.Price, o.Amount)
}

func (o *Order) Renew() {
	o.ID = uuid.New()
	o.Status = OsNew
}

package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

const OrderCollectionName = "orders"

type Order struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Side         Side               `json:"side" bson:"side"`
	Symbol       Symbol             `json:"symbol" bson:"symbol"`
	TotalAmount  float64            `json:"totalAmount" bson:"total_amount"`
	RemainAmount float64            `json:"remainAmount" bson:"remain_amount"`
	Price        uint32             `json:"price" bson:"price"`
	Next         *Order             `json:"next" bson:"next"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updated_at"`
}

func (o *Order) String() string {
	return fmt.Sprintf("order{id:%v,isBuy:%v,price:%v,total:%v,remain:%v}", o.ID, o.Side, o.Price, o.TotalAmount, o.RemainAmount)
}

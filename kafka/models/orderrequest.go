package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"orderbook/core/models"
)

type OrderRequest struct {
	Side   models.Side   `json:"side"`
	Symbol models.Symbol `json:"symbol"`
	Amount float64       `json:"amount"`
	Price  uint32        `json:"price"`
}

func NewOrder(or OrderRequest) *models.Order {
	return &models.Order{
		ID:           primitive.NewObjectID(),
		Side:         or.Side,
		Symbol:       or.Symbol,
		TotalAmount:  or.Amount,
		RemainAmount: or.Amount,
		Price:        or.Price,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

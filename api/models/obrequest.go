package models

import "orderbook/core/models"

type OrderbookRequest struct {
	Symbol models.Symbol `json:"symbol" binding:"required"`
	Limit  int           `json:"limit"`
}

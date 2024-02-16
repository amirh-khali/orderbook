package models

type OrderbookRequest struct {
	Symbol string `json:"symbol" binding:"required"`
	Limit  int    `json:"limit"`
}

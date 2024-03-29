package models

type OrderbookResponse struct {
	Bids   []OrderRecord `json:"bids"`
	Asks   []OrderRecord `json:"asks"`
	MinAsk uint32        `json:"minAsk"`
	MaxBid uint32        `json:"maxBid"`
}

type OrderRecord [2]string

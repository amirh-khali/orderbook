package models

type OrderbookResponse struct {
	LastUpdateID int           `json:"lastUpdateID"`
	Bids         []OrderRecord `json:"bids"`
	Asks         []OrderRecord `json:"asks"`
}

type OrderRecord [2]string

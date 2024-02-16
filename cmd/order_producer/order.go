package main

type Order struct {
	Side   Side    `json:"side"`
	Symbol Symbol  `json:"symbol"`
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type Symbol int

const (
	BTCUSDT Symbol = iota
	ETHUSDT
	BTCIRT
	ETHIRT
)

type Side int

const (
	Buy Side = iota
	Sell
)

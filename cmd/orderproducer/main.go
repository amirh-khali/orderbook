package main

import (
	"math/rand"
	"time"

	"orderbook/env"
	"orderbook/kafka"
	"orderbook/scheduler"
)

const (
	minAmount = 100
	maxAmount = 100000
)

var minPrice map[Symbol]int
var maxPrice map[Symbol]int
var priceMultiplier map[Symbol]int
var symbols []Symbol
var sides []Side

func produce() {
	sd := sides[rand.Intn(len(sides))]
	sb := symbols[rand.Intn(len(symbols))]
	a := float64(minAmount+rand.Intn(maxAmount-minAmount)) / 100000
	p := float64(priceMultiplier[sb]) * float64(minPrice[sb]+rand.Intn(maxPrice[sb]-minPrice[sb]))

	kafka.Produce(Order{Side: sd, Symbol: sb, Amount: a, Price: p}, env.ENV.KafkaTopic)
}

func constantsInit() {
	minPrice = map[Symbol]int{
		BTCUSDT: 490,
		ETHUSDT: 2780,
		BTCIRT:  28950,
		ETHIRT:  1550,
	}
	maxPrice = map[Symbol]int{
		BTCUSDT: 510,
		ETHUSDT: 2800,
		BTCIRT:  28970,
		ETHIRT:  1560,
	}
	priceMultiplier = map[Symbol]int{
		BTCUSDT: 100,
		ETHUSDT: 1,
		BTCIRT:  100000,
		ETHIRT:  100000,
	}
	symbols = []Symbol{BTCUSDT, ETHUSDT, BTCIRT, ETHIRT}
	sides = []Side{Buy, Sell}
}

func main() {
	env.LoadEnv()
	constantsInit()

	kafka.InitProducer()

	scheduler.Init()
	scheduler.AddJob(time.Second, produce)
	scheduler.Start()
}

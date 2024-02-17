package main

import (
	"math/rand"
	"time"

	"orderbook/db/models"
	"orderbook/env"
	"orderbook/kafka"
	"orderbook/scheduler"
)

const (
	minAmount = 100
	maxAmount = 100000
)

var minPrice map[models.Symbol]int
var maxPrice map[models.Symbol]int
var priceMultiplier map[models.Symbol]int
var symbols []models.Symbol
var sides []models.Side

func produce() {
	sd := sides[rand.Intn(len(sides))]
	sb := symbols[rand.Intn(len(symbols))]
	a := float64(minAmount+rand.Intn(maxAmount-minAmount)) / 100000
	p := uint32(priceMultiplier[sb]*minPrice[sb] + rand.Intn(maxPrice[sb]-minPrice[sb]))

	kafka.Produce(models.Order{Side: sd, Symbol: sb, Amount: a, Price: p}, env.ENV.KafkaTopic)
}

func constantsInit() {
	minPrice = map[models.Symbol]int{
		models.BTCUSDT: 490,
		models.ETHUSDT: 2780,
		models.BTCIRT:  28950,
		models.ETHIRT:  1550,
	}
	maxPrice = map[models.Symbol]int{
		models.BTCUSDT: 510,
		models.ETHUSDT: 2800,
		models.BTCIRT:  28970,
		models.ETHIRT:  1560,
	}
	priceMultiplier = map[models.Symbol]int{
		models.BTCUSDT: 10,
		models.ETHUSDT: 1,
		models.BTCIRT:  1000,
		models.ETHIRT:  1000,
	}
	symbols = []models.Symbol{models.BTCUSDT, models.ETHUSDT, models.BTCIRT, models.ETHIRT}
	sides = []models.Side{models.Buy, models.Sell}
}

func main() {
	env.LoadEnv()
	constantsInit()

	kafka.InitProducer(env.ENV.BootstrapServers)

	scheduler.Init()
	scheduler.AddJob(time.Second, produce)
	scheduler.Start()
}

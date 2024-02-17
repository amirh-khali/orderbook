package main

import (
	"github.com/gin-gonic/gin"
	"orderbook/api"
	"orderbook/db"
	"orderbook/env"
	"orderbook/kafka"
)

func newRouter() *gin.Engine {
	r := gin.Default()

	h := api.NewOrderbookHandler()
	r.GET("/orderbook", h.Get)

	return r
}

func main() {
	env.LoadEnv()
	r := newRouter()

	obMap := db.NewOrderBookMap()

	kafka.InitConsumer(env.ENV.BootstrapServers, env.ENV.KafkaGroupID, env.ENV.KafkaAutoOffsetReset)
	kafka.Subscribe(env.ENV.KafkaTopic)
	go kafka.StartConsume(obMap)

	err := r.Run()
	if err != nil {
		return
	}
}

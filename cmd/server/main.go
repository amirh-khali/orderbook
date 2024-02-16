package main

import (
	"github.com/gin-gonic/gin"
	"orderbook/api"
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

	kafka.InitConsumer(env.ENV.BootstrapServers, env.ENV.KafkaGroupID, env.ENV.KafkaAutoOffsetReset)
	kafka.Subscribe(env.ENV.KafkaTopic)
	kafka.StartConsume()

	err := r.Run()
	if err != nil {
		return
	}
}

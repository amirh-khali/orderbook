package main

import (
	"github.com/gin-gonic/gin"
	"orderbook/api"
	"orderbook/core"
	"orderbook/db/mongo"
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

	mongo.InitClient(env.ENV.MongoURI, env.ENV.MongoDatabaseName)
	mongo.RegisterCollections()

	core.Init()

	kafka.InitConsumer(env.ENV.BootstrapServers, env.ENV.KafkaGroupID, env.ENV.KafkaAutoOffsetReset)
	kafka.Subscribe(env.ENV.KafkaTopic)
	go kafka.StartConsume()

	err := r.Run()
	if err != nil {
		return
	}
}

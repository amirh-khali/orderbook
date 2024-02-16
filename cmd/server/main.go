package main

import (
	"github.com/gin-gonic/gin"
	"orderbook/api"
)

func main() {
	r := gin.Default()

	h := api.NewOrderbookHandler()
	r.GET("/orderbook", h.Get)

	err := r.Run()
	if err != nil {
		return
	}
}

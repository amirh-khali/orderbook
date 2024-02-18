package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"orderbook/api/models"
	"orderbook/core"
)

type OrderbookHandler struct{}

func NewOrderbookHandler() *OrderbookHandler {
	return &OrderbookHandler{}
}

func (h OrderbookHandler) Get(c *gin.Context) {
	var obRequest models.OrderbookRequest
	if err := c.ShouldBindJSON(&obRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if obRequest.Limit == 0 {
		obRequest.Limit = 10
	}

	ob := (*core.OrderbookMap)[obRequest.Symbol]

	c.JSON(http.StatusOK, gin.H{
		"data": models.OrderbookResponse{
			Bids:   ob.GetBidsTable(obRequest.Limit),
			Asks:   ob.GetAsksTable(obRequest.Limit),
			MinAsk: ob.MinAsk,
			MaxBid: ob.MaxBid,
		},
	})
}

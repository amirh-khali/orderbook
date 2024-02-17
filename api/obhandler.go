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

	ob := (*core.OrderbookMap)[obRequest.Symbol]

	c.JSON(http.StatusOK, gin.H{
		"data": models.OrderbookResponse{
			LastUpdateID: 1027024,
			Bids:         ob.GetBidsTable(10),
			Asks:         ob.GetAsksTable(10),
			MinAsk:       ob.MinAsk,
			MaxBid:       ob.MaxBid,
		},
	})
}

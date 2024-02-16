package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"orderbook/api/models"
)

type OrderbookHandler struct{}

func NewOrderbookHandler() *OrderbookHandler {
	return &OrderbookHandler{}
}

func (h OrderbookHandler) Get(c *gin.Context) {
	ob := models.OrderbookResponse{
		LastUpdateID: 1027024,
		Bids:         []models.OrderRecord{{"4.00000000", "431.00000000"}},
		Asks:         []models.OrderRecord{{"4.00000000", "12.00000000"}},
	}

	c.JSON(http.StatusOK, gin.H{"data": ob})
}

package handlers

import (
	"net/http"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/internal/models"
	"github.com/4planet/backend/pkg/prices"
	"github.com/gin-gonic/gin"
)

type PricesHandler struct {
	pricesService *prices.Service
	config        *config.Config
}

func NewPricesHandler(pricesService *prices.Service, config *config.Config) *PricesHandler {
	return &PricesHandler{
		pricesService: pricesService,
		config:        config,
	}
}

// GetPrices retrieves all tree prices
func (h *PricesHandler) GetPrices(c *gin.Context) {
	prices, err := h.pricesService.GetPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prices"})
		return
	}

	c.JSON(http.StatusOK, prices)
}

// GetPriceByCurrency retrieves a tree price by currency
func (h *PricesHandler) GetPriceByCurrency(c *gin.Context) {
	currencyStr := c.Param("currency")

	// Parse currency
	currency := models.Currency(currencyStr)

	// Validate currency (optional, since the database will handle invalid currencies)

	price, err := h.pricesService.GetPriceByCurrency(currency)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Price not found for currency"})
		return
	}

	c.JSON(http.StatusOK, price)
}

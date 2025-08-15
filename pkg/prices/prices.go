package prices

import (
	"github.com/4planet/backend/internal/database"
	"github.com/4planet/backend/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService() *Service {
	return &Service{
		db: database.GetDB(),
	}
}

// GetPrices retrieves all tree prices as a map of currency to price
func (s *Service) GetPrices() (map[string]int64, error) {
	var prices []models.TreePrice
	err := s.db.Order("currency").Find(&prices).Error
	if err != nil {
		return nil, err
	}

	// Convert to map for easier frontend consumption
	pricesMap := make(map[string]int64)
	for _, price := range prices {
		pricesMap[string(price.Currency)] = price.PriceMinor
	}

	return pricesMap, nil
}

// GetPriceByCurrency retrieves a tree price by currency
func (s *Service) GetPriceByCurrency(currency models.Currency) (*models.TreePrice, error) {
	var price models.TreePrice
	err := s.db.Where("currency = ?", currency).First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

// UpdatePrice updates or creates a tree price for a specific currency
func (s *Service) UpdatePrice(currency models.Currency, priceMinor int64) error {
	price := models.TreePrice{
		Currency:   currency,
		PriceMinor: priceMinor,
	}

	return s.db.Save(&price).Error
}

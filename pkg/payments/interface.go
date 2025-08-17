package payments

import (
	"errors"
)

// Payment provider errors
var (
	ErrProviderNotSupported = errors.New("payment provider not supported")
	ErrProviderDisabled     = errors.New("payment provider is disabled")
)

// PaymentService defines the interface that any payment provider must implement
type PaymentService interface {
	// CreatePaymentIntent creates a payment intent for one-time payment
	CreatePaymentIntent(req *PaymentIntentRequest, authUserID string) (*PaymentIntentResponse, error)

	// CreateSubscriptionIntent creates a subscription intent for recurring payments
	CreateSubscriptionIntent(req *SubscriptionIntentRequest, authUserID string) (*SubscriptionIntentResponse, error)

	// ProcessWebhook processes webhook events from the payment provider
	ProcessWebhook(payload []byte, signature string) error

	// GetProviderName returns the name of the payment provider
	GetProviderName() string
}

// PaymentProviderFactory creates payment services based on provider name
type PaymentProviderFactory struct {
	configs map[string]PaymentProviderConfig
}

// PaymentProviderConfig holds configuration for a specific payment provider
type PaymentProviderConfig struct {
	ProviderName string
	PublicID     string
	Secret       string
	BaseURL      string
	Enabled      bool
}

// NewPaymentProviderFactory creates a new factory with provider configurations
func NewPaymentProviderFactory(configs map[string]PaymentProviderConfig) *PaymentProviderFactory {
	return &PaymentProviderFactory{
		configs: configs,
	}
}

// CreateService creates a payment service for the specified provider
func (f *PaymentProviderFactory) CreateService(providerName string) (PaymentService, error) {
	config, exists := f.configs[providerName]
	if !exists {
		return nil, ErrProviderNotSupported
	}

	if !config.Enabled {
		return nil, ErrProviderDisabled
	}

	switch providerName {
	case "cloudpayments":
		return NewCloudPaymentsService(config.PublicID, config.Secret, config.BaseURL), nil
	// Add more providers here as they become available
	// case "stripe":
	//     return NewStripeService(config.PublicID, config.Secret, config.BaseURL), nil
	// case "paypal":
	//     return NewPayPalService(config.PublicID, config.Secret, config.BaseURL), nil
	default:
		return nil, ErrProviderNotSupported
	}
}

// GetSupportedProviders returns a list of supported payment providers
func (f *PaymentProviderFactory) GetSupportedProviders() []string {
	var providers []string
	for name, config := range f.configs {
		if config.Enabled {
			providers = append(providers, name)
		}
	}
	return providers
}

// IsProviderSupported checks if a provider is supported and enabled
func (f *PaymentProviderFactory) IsProviderSupported(providerName string) bool {
	config, exists := f.configs[providerName]
	return exists && config.Enabled
}

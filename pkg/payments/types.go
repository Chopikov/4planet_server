package payments

// SubscriptionDescriptions holds subscription description configurations
type SubscriptionDescriptions struct {
	Monthly string
	Yearly  string
	Custom  string
}

// PaymentTexts holds all payment-related text configurations
type PaymentTexts struct {
	DefaultDonationDescription  string
	BaseSubscriptionDescription string
	SubscriptionDescriptions    SubscriptionDescriptions
}

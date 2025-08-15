package models

import (
	"database/sql/driver"
	"fmt"
)

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusPlanned    ProjectStatus = "planned"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusCompleted  ProjectStatus = "completed"
)

func (ps ProjectStatus) String() string {
	return string(ps)
}

func (ps *ProjectStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*ps = ProjectStatus(v)
	case []byte:
		*ps = ProjectStatus(string(v))
	default:
		return fmt.Errorf("cannot scan %T into ProjectStatus", value)
	}
	return nil
}

func (ps ProjectStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// NewsType represents the type of news item
type NewsType string

const (
	NewsTypeAchievement NewsType = "achievement"
	NewsTypeInvite      NewsType = "invite"
	NewsTypeUpdate      NewsType = "update"
)

// IsValid checks if the NewsType value is valid
func (nt NewsType) IsValid() bool {
	switch nt {
	case NewsTypeAchievement, NewsTypeInvite, NewsTypeUpdate:
		return true
	default:
		return false
	}
}

// ParseNewsType safely parses a string to NewsType, returning false if invalid
func ParseNewsType(s string) (NewsType, bool) {
	nt := NewsType(s)
	return nt, nt.IsValid()
}

func (nt NewsType) String() string {
	return string(nt)
}

func (nt *NewsType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*nt = NewsType(v)
	case []byte:
		*nt = NewsType(string(v))
	default:
		return fmt.Errorf("cannot scan %T into NewsType", value)
	}
	return nil
}

func (nt NewsType) Value() (driver.Value, error) {
	return string(nt), nil
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCanceled  PaymentStatus = "canceled"
)

func (ps PaymentStatus) String() string {
	return string(ps)
}

func (ps *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*ps = PaymentStatus(v)
	case []byte:
		*ps = PaymentStatus(string(v))
	default:
		return fmt.Errorf("cannot scan %T into PaymentStatus", value)
	}
	return nil
}

func (ps PaymentStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// PaymentProvider represents the payment provider
type PaymentProvider string

const (
	PaymentProviderCloudPayments PaymentProvider = "cloudpayments"
	PaymentProviderKaspi         PaymentProvider = "kaspi"
	PaymentProviderPayPal        PaymentProvider = "paypal"
	PaymentProviderTribute       PaymentProvider = "tribute"
)

func (pp PaymentProvider) String() string {
	return string(pp)
}

func (pp *PaymentProvider) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*pp = PaymentProvider(v)
	case []byte:
		*pp = PaymentProvider(string(v))
	default:
		return fmt.Errorf("cannot scan %T into PaymentProvider", value)
	}
	return nil
}

func (pp PaymentProvider) Value() (driver.Value, error) {
	return string(pp), nil
}

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive     SubscriptionStatus = "active"
	SubscriptionStatusPastDue    SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled   SubscriptionStatus = "canceled"
	SubscriptionStatusPaused     SubscriptionStatus = "paused"
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete"
)

func (ss SubscriptionStatus) String() string {
	return string(ss)
}

func (ss *SubscriptionStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*ss = SubscriptionStatus(v)
	case []byte:
		*ss = SubscriptionStatus(string(v))
	default:
		return fmt.Errorf("cannot scan %T into SubscriptionStatus", value)
	}
	return nil
}

func (ss SubscriptionStatus) Value() (driver.Value, error) {
	return string(ss), nil
}

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusPending UserStatus = "pending"
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
)

func (us UserStatus) String() string {
	return string(us)
}

func (us *UserStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*us = UserStatus(v)
	case []byte:
		*us = UserStatus(string(v))
	default:
		return fmt.Errorf("cannot scan %T into UserStatus", value)
	}
	return nil
}

func (us UserStatus) Value() (driver.Value, error) {
	return string(us), nil
}

// MediaKind represents the kind of media file
type MediaKind string

const (
	MediaKindImage    MediaKind = "image"
	MediaKindVideo    MediaKind = "video"
	MediaKindDocument MediaKind = "document"
)

func (mk MediaKind) String() string {
	return string(mk)
}

func (mk *MediaKind) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*mk = MediaKind(v)
	case []byte:
		*mk = MediaKind(string(v))
	default:
		return fmt.Errorf("cannot scan %T into MediaKind", value)
	}
	return nil
}

func (mk MediaKind) Value() (driver.Value, error) {
	return string(mk), nil
}

// Currency represents the currency code
type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyKZT Currency = "KZT"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

func (c Currency) String() string {
	return string(c)
}

func (c *Currency) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*c = Currency(v)
	case []byte:
		*c = Currency(string(v))
	default:
		return fmt.Errorf("cannot scan %T into Currency", value)
	}
	return nil
}

func (c Currency) Value() (driver.Value, error) {
	return string(c), nil
}

// ShareKind represents the kind of share link
type ShareKind string

const (
	ShareKindProfile  ShareKind = "profile"
	ShareKindDonation ShareKind = "donation"
)

func (sk ShareKind) String() string {
	return string(sk)
}

func (sk *ShareKind) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*sk = ShareKind(v)
	case []byte:
		*sk = ShareKind(string(v))
	default:
		return fmt.Errorf("cannot scan %T into ShareKind", value)
	}
	return nil
}

func (sk ShareKind) Value() (driver.Value, error) {
	return string(sk), nil
}

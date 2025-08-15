package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the users table (profile data)
type User struct {
	ID             uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID     string     `gorm:"column:auth_user_id;uniqueIndex;type:text;not null"`
	Username       *string    `gorm:"column:username;uniqueIndex;type:text"`
	DisplayName    *string    `gorm:"column:display_name;type:text"`
	AvatarURL      *string    `gorm:"column:avatar_url;type:text"`
	Email          string     `gorm:"column:email;type:text;not null"`
	TotalTrees     int        `gorm:"column:total_trees;type:int;not null;default:0"`
	DonationsCount int        `gorm:"column:donations_count;type:int;not null;default:0"`
	LastDonationAt *time.Time `gorm:"column:last_donation_at;type:timestamptz"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	// Relationships
	Sessions         []Session         `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	Subscriptions    []Subscription    `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	Payments         []Payment         `gorm:"foreignKey:AuthUserID;constraint:OnDelete:SET NULL"`
	Donations        []Donation        `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	UserAchievements []UserAchievement `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	ShareTokens      []ShareToken      `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
}

// UserAuth represents user authentication data
type UserAuth struct {
	ID           uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID   string     `gorm:"column:auth_user_id;uniqueIndex;type:text;not null"`
	Email        string     `gorm:"column:email;uniqueIndex;type:text;not null"`
	PasswordHash *string    `gorm:"column:password_hash;type:text"`
	Status       UserStatus `gorm:"column:status;type:user_status;not null;default:'pending'"`
	VerifiedAt   *time.Time `gorm:"column:verified_at;type:timestamptz"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`

	// Relationships
	EmailVerificationTokens []EmailVerificationToken `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
	PasswordResetTokens     []PasswordResetToken     `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}

func (UserAuth) TableName() string {
	return "user_auth"
}

// Session represents the sessions table
type Session struct {
	ID         uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID string     `gorm:"column:auth_user_id;type:text;not null;index"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;type:timestamptz;not null"`
	RevokedAt  *time.Time `gorm:"column:revoked_at;type:timestamptz"`
	UserAgent  *string    `gorm:"column:user_agent;type:text"`
	IPAddr     *string    `gorm:"column:ip_addr;type:inet"`

	// Relationships
	User User `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
}

func (Session) TableName() string {
	return "sessions"
}

// EmailVerificationToken represents the email_verification_tokens table
type EmailVerificationToken struct {
	ID         uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID string     `gorm:"column:auth_user_id;type:text;not null;index"`
	Token      string     `gorm:"column:token;type:text;uniqueIndex;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;type:timestamptz;not null"`
	UsedAt     *time.Time `gorm:"column:used_at;type:timestamptz"`

	// Relationships
	User User `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// PasswordResetToken represents the password_reset_tokens table
type PasswordResetToken struct {
	ID         uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID string     `gorm:"column:auth_user_id;type:text;not null;index"`
	Token      string     `gorm:"column:token;type:text;uniqueIndex;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	ExpiresAt  time.Time  `gorm:"column:expires_at;type:timestamptz;not null"`
	UsedAt     *time.Time `gorm:"column:used_at;type:timestamptz"`

	// Relationships
	User User `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// TreePrice represents the tree_prices table
type TreePrice struct {
	Currency   Currency  `gorm:"column:currency;primaryKey;type:text"`
	PriceMinor int64     `gorm:"column:price_minor;type:bigint;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (TreePrice) TableName() string {
	return "tree_prices"
}

// Project represents the projects table
type Project struct {
	ID              uuid.UUID     `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string        `gorm:"column:title;type:text;not null"`
	Description     *string       `gorm:"column:description;type:text"`
	Status          ProjectStatus `gorm:"column:status;type:project_status;not null;default:'planned';index"`
	StartsAt        *time.Time    `gorm:"column:starts_at;type:timestamptz"`
	EndsAt          *time.Time    `gorm:"column:ends_at;type:timestamptz"`
	CountryCode     *string       `gorm:"column:country_code;type:text"`
	Region          *string       `gorm:"column:region;type:text"`
	LocationGeoJSON interface{}   `gorm:"column:location_geojson;type:jsonb;not null;index:type:gin"`
	TreesTarget     *int          `gorm:"column:trees_target;type:integer"`
	TreesPlanted    *int          `gorm:"column:trees_planted;type:integer"`
	CoverURL        *string       `gorm:"column:cover_url;type:text"`
	CreatedAt       time.Time     `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	// Relationships
	MediaFiles []MediaFile `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	News       []News      `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	Donations  []Donation  `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
}

func (Project) TableName() string {
	return "projects"
}

// MediaFile represents the media_files table
type MediaFile struct {
	ID        uuid.UUID   `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	ProjectID uuid.UUID   `gorm:"column:project_id;type:uuid;not null;index"`
	Kind      MediaKind   `gorm:"column:kind;type:media_kind;not null;default:'image'"`
	URL       string      `gorm:"column:url;type:text;not null"`
	MimeType  *string     `gorm:"column:mime_type;type:text"`
	Title     *string     `gorm:"column:title;type:text"`
	AltText   *string     `gorm:"column:alt_text;type:text"`
	Meta      interface{} `gorm:"column:meta;type:jsonb;default:'{}'::jsonb"`
	CreatedAt time.Time   `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	// Relationships
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

func (MediaFile) TableName() string {
	return "media_files"
}

// News represents the news table
type News struct {
	ID          uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Type        NewsType   `gorm:"column:type;type:news_type;not null"`
	Title       string     `gorm:"column:title;type:text;not null"`
	BodyMD      *string    `gorm:"column:body_md;type:text"`
	CoverURL    *string    `gorm:"column:cover_url;type:text"`
	ProjectID   *uuid.UUID `gorm:"column:project_id;type:uuid;index"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	PublishedAt *time.Time `gorm:"column:published_at;type:timestamptz;index"`

	// Relationships
	Project *Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
}

func (News) TableName() string {
	return "news"
}

// Achievement represents the achievements table
type Achievement struct {
	ID             uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Code           string    `gorm:"column:code;type:text;uniqueIndex;not null"`
	Title          string    `gorm:"column:title;type:text;not null"`
	Description    *string   `gorm:"column:description;type:text"`
	ThresholdTrees *int      `gorm:"column:threshold_trees;type:integer"`
	ImageURL       *string   `gorm:"column:image_url;type:text"`

	// Relationships
	UserAchievements []UserAchievement `gorm:"foreignKey:AchievementID;constraint:OnDelete:CASCADE"`
}

func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement represents the user_achievements table
type UserAchievement struct {
	AuthUserID    string    `gorm:"column:auth_user_id;type:text;primaryKey"`
	AchievementID uuid.UUID `gorm:"column:achievement_id;type:uuid;primaryKey"`
	AwardedAt     time.Time `gorm:"column:awarded_at;type:timestamptz;not null;default:now()"`
	Reason        *string   `gorm:"column:reason;type:text"`

	// Relationships
	User        User        `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE" json:"-"`
	Achievement Achievement `gorm:"foreignKey:AchievementID;constraint:OnDelete:CASCADE"`
}

func (UserAchievement) TableName() string {
	return "user_achievements"
}

// Subscription represents the subscriptions table
type Subscription struct {
	ID                     uuid.UUID          `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID             string             `gorm:"column:auth_user_id;type:text;not null;index"`
	Provider               PaymentProvider    `gorm:"column:provider;type:payment_provider;not null"`
	ProviderCustomerID     *string            `gorm:"column:provider_customer_id;type:text"`
	ProviderSubscriptionID *string            `gorm:"column:provider_subscription_id;type:text;uniqueIndex"`
	AmountMinor            int64              `gorm:"column:amount_minor;type:bigint;not null"`
	Currency               Currency           `gorm:"column:currency;type:text;not null"`
	IntervalMonths         int                `gorm:"column:interval_months;type:integer;not null;default:1"`
	Status                 SubscriptionStatus `gorm:"column:status;type:subscription_status;not null"`
	StartedAt              time.Time          `gorm:"column:started_at;type:timestamptz;not null;default:now()"`
	CanceledAt             *time.Time         `gorm:"column:canceled_at;type:timestamptz"`
	Meta                   interface{}        `gorm:"column:meta;type:jsonb;default:'{}'::jsonb"`

	// Relationships
	User     User      `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE" json:"-"`
	Payments []Payment `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:SET NULL"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

// Payment represents the payments table
type Payment struct {
	ID                uuid.UUID       `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Provider          PaymentProvider `gorm:"column:provider;type:payment_provider;not null"`
	ProviderPaymentID *string         `gorm:"column:provider_payment_id;type:text;uniqueIndex"`
	AuthUserID        *string         `gorm:"column:auth_user_id;type:text;index"`
	SubscriptionID    *uuid.UUID      `gorm:"column:subscription_id;type:uuid;index"`
	AmountMinor       int64           `gorm:"column:amount_minor;type:bigint;not null"`
	Currency          Currency        `gorm:"column:currency;type:text;not null"`
	Status            PaymentStatus   `gorm:"column:status;type:payment_status;not null;index"`
	OccurredAt        *time.Time      `gorm:"column:occurred_at;type:timestamptz"`
	Meta              interface{}     `gorm:"column:meta;type:jsonb;default:'{}'::jsonb"`
	CreatedAt         time.Time       `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	// Relationships
	User         *User         `gorm:"foreignKey:AuthUserID;constraint:OnDelete:SET NULL"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:SET NULL"`
	Donation     *Donation     `gorm:"foreignKey:PaymentID;constraint:OnDelete:RESTRICT"`
}

func (Payment) TableName() string {
	return "payments"
}

// Donation represents the donations table
type Donation struct {
	ID             uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID     string     `gorm:"column:auth_user_id;type:text;not null;index"`
	PaymentID      uuid.UUID  `gorm:"column:payment_id;type:uuid;uniqueIndex;not null"`
	ProjectID      *uuid.UUID `gorm:"column:project_id;type:uuid;index"`
	ReferralUserID *string    `gorm:"column:referral_user_id;type:text;index"`
	TreesCount     int        `gorm:"column:trees_count;type:integer;not null"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	// Relationships
	User         User         `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE" json:"-"`
	Payment      Payment      `gorm:"foreignKey:PaymentID;constraint:OnDelete:RESTRICT"`
	Project      *Project     `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL" json:"-"`
	ReferralUser *User        `gorm:"foreignKey:ReferralUserID;constraint:OnDelete:SET NULL" json:"-"`
	ShareTokens  []ShareToken `gorm:"foreignKey:RefID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Donation) TableName() string {
	return "donations"
}

// ShareToken represents the share_tokens table
type ShareToken struct {
	ID         uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthUserID string     `gorm:"column:auth_user_id;type:text;not null;index"`
	Kind       ShareKind  `gorm:"column:kind;type:text;not null"`
	RefID      *uuid.UUID `gorm:"column:ref_id;type:uuid"`
	Slug       string     `gorm:"column:slug;type:text;uniqueIndex;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`

	// Relationships
	User User `gorm:"foreignKey:AuthUserID;constraint:OnDelete:CASCADE"`
}

func (ShareToken) TableName() string {
	return "share_tokens"
}

// WebhookEvent represents the webhook_events table
type WebhookEvent struct {
	ID               uuid.UUID       `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Provider         PaymentProvider `gorm:"column:provider;type:payment_provider;not null;index"`
	EventType        string          `gorm:"column:event_type;type:text;not null"`
	EventIdempotency *string         `gorm:"column:event_idempotency;type:text;uniqueIndex"`
	ReceivedAt       time.Time       `gorm:"column:received_at;type:timestamptz;not null;default:now();index"`
	RawPayload       interface{}     `gorm:"column:raw_payload;type:jsonb;not null"`
	SignatureOK      bool            `gorm:"column:signature_ok;type:boolean;not null"`
	ProcessedOK      bool            `gorm:"column:processed_ok;type:boolean;not null;default:false"`
	ProcessingError  *string         `gorm:"column:processing_error;type:text"`
}

func (WebhookEvent) TableName() string {
	return "webhook_events"
}

// UserStats represents the user_stats view
type UserStats struct {
	AuthUserID     string     `gorm:"column:auth_user_id"`
	TotalTrees     int        `gorm:"column:total_trees"`
	DonationsCount int        `gorm:"column:donations_count"`
	LastDonationAt *time.Time `gorm:"column:last_donation_at"`
}

func (UserStats) TableName() string {
	return "user_stats"
}

-- Initial database schema for 4Planet backend
-- Generated from GORM models to ensure perfect alignment

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =============================
-- ENUM TYPES
-- =============================

CREATE TYPE project_status AS ENUM ('planned', 'in_progress', 'completed');
CREATE TYPE news_type AS ENUM ('achievement', 'invite', 'update');
CREATE TYPE payment_status AS ENUM ('pending', 'succeeded', 'failed', 'refunded', 'canceled');
CREATE TYPE payment_provider AS ENUM ('cloudpayments', 'kaspi', 'paypal', 'tribute');
CREATE TYPE subscription_status AS ENUM ('active', 'past_due', 'canceled', 'paused', 'incomplete');
CREATE TYPE user_status AS ENUM ('pending', 'active', 'blocked');
CREATE TYPE media_kind AS ENUM ('image', 'video', 'document');

-- =============================
-- CORE TABLES
-- =============================

-- UserAuth table (authentication data)
CREATE TABLE user_auth (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text UNIQUE NOT NULL,
    email text UNIQUE NOT NULL,
    password_hash text,
    status user_status NOT NULL DEFAULT 'pending',
    verified_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Users table (profile data)
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text UNIQUE NOT NULL,
    username text UNIQUE,
    display_name text,
    avatar_url text,
    email text NOT NULL,
    total_trees integer NOT NULL DEFAULT 0,
    donations_count integer NOT NULL DEFAULT 0,
    last_donation_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_users_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE
);

-- Sessions table
CREATE TABLE sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    user_agent text,
    ip_addr inet,
    CONSTRAINT fk_sessions_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE
);

-- Email verification tokens table
CREATE TABLE email_verification_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    token text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    CONSTRAINT fk_email_verif_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE
);

-- Password reset tokens table
CREATE TABLE password_reset_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    token text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    CONSTRAINT fk_password_reset_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE
);

-- Tree prices table
CREATE TABLE tree_prices (
    currency text PRIMARY KEY,
    price_minor bigint NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Projects table
CREATE TABLE projects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description text,
    status project_status NOT NULL DEFAULT 'planned',
    starts_at timestamptz,
    ends_at timestamptz,
    country_code text,
    region text,
    location_geojson jsonb NOT NULL,
    trees_target integer,
    trees_planted integer,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- Media files table
CREATE TABLE media_files (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL,
    kind media_kind NOT NULL DEFAULT 'image',
    url text NOT NULL,
    mime_type text,
    title text,
    alt_text text,
    meta jsonb DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_media_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- News table
CREATE TABLE news (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type news_type NOT NULL,
    title text NOT NULL,
    body_md text,
    cover_url text,
    project_id uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    published_at timestamptz,
    CONSTRAINT fk_news_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);

-- Achievements table
CREATE TABLE achievements (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text UNIQUE NOT NULL,
    title text NOT NULL,
    description text,
    threshold_trees integer,
    image_url text
);

-- User achievements table
CREATE TABLE user_achievements (
    auth_user_id text NOT NULL,
    achievement_id uuid NOT NULL,
    awarded_at timestamptz NOT NULL DEFAULT now(),
    reason text,
    PRIMARY KEY (auth_user_id, achievement_id),
    CONSTRAINT fk_user_achievements_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_achievements_achievement FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE
);

-- Subscriptions table
CREATE TABLE subscriptions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    provider payment_provider NOT NULL,
    provider_customer_id text,
    provider_subscription_id text UNIQUE,
    amount_minor bigint NOT NULL,
    currency text NOT NULL,
    interval_months integer NOT NULL DEFAULT 1,
    status subscription_status NOT NULL,
    started_at timestamptz NOT NULL DEFAULT now(),
    canceled_at timestamptz,
    meta jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT fk_subscriptions_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE
);

-- Payments table
CREATE TABLE payments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    provider payment_provider NOT NULL,
    provider_payment_id text UNIQUE,
    auth_user_id text,
    subscription_id uuid,
    amount_minor bigint NOT NULL,
    currency text NOT NULL,
    status payment_status NOT NULL,
    occurred_at timestamptz,
    meta jsonb DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_payments_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE SET NULL,
    CONSTRAINT fk_payments_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE SET NULL
);

-- Donations table
CREATE TABLE donations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    payment_id uuid UNIQUE NOT NULL,
    project_id uuid,
    trees_count integer NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_donations_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE,
    CONSTRAINT fk_donations_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE RESTRICT,
    CONSTRAINT fk_donations_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);

-- Share tokens table
CREATE TABLE share_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    kind text NOT NULL,
    ref_id uuid,
    slug text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_share_tokens_user_auth FOREIGN KEY (auth_user_id) REFERENCES user_auth(auth_user_id) ON DELETE CASCADE
);

-- Webhook events table
CREATE TABLE webhook_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    provider payment_provider NOT NULL,
    event_type text NOT NULL,
    event_idempotency text UNIQUE,
    received_at timestamptz NOT NULL DEFAULT now(),
    raw_payload jsonb NOT NULL,
    signature_ok boolean NOT NULL,
    processed_ok boolean NOT NULL DEFAULT false,
    processing_error text
);

-- =============================
-- INDEXES
-- =============================

-- UserAuth indexes
CREATE INDEX idx_user_auth_email ON user_auth(email);
CREATE INDEX idx_user_auth_status ON user_auth(status);

-- Users indexes
CREATE INDEX idx_users_auth_user_id ON users(auth_user_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Sessions indexes
CREATE INDEX idx_sessions_user ON sessions(auth_user_id);
CREATE INDEX idx_sessions_valid ON sessions(expires_at) WHERE revoked_at IS NULL;

-- Email verification tokens indexes
CREATE INDEX idx_email_verif_user ON email_verification_tokens(auth_user_id);

-- Password reset tokens indexes
CREATE INDEX idx_password_reset_user ON password_reset_tokens(auth_user_id);

-- Projects indexes
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_geo ON projects USING GIN (location_geojson);

-- Media files indexes
CREATE INDEX idx_media_project ON media_files(project_id);

-- News indexes
CREATE INDEX idx_news_published ON news(published_at DESC);

-- Subscriptions indexes
CREATE INDEX idx_subscriptions_user ON subscriptions(auth_user_id);

-- Payments indexes
CREATE INDEX idx_payments_user ON payments(auth_user_id, created_at DESC);
CREATE INDEX idx_payments_status ON payments(status);

-- Donations indexes
CREATE INDEX idx_donations_user_created ON donations(auth_user_id, created_at DESC);
CREATE INDEX idx_donations_project ON donations(project_id);

-- Share tokens indexes
CREATE INDEX idx_share_tokens_user ON share_tokens(auth_user_id);

-- Webhook events indexes
CREATE INDEX idx_webhook_events_provider ON webhook_events(provider, received_at DESC);

-- =============================
-- VIEWS
-- =============================

-- User stats view for consistency checking
CREATE OR REPLACE VIEW user_stats AS
SELECT
    u.auth_user_id,
    COALESCE(SUM(d.trees_count), 0)::int AS total_trees,
    COUNT(d.*)::int AS donations_count,
    MAX(d.created_at) AS last_donation_at
FROM users u
LEFT JOIN donations d ON d.auth_user_id = u.auth_user_id
GROUP BY u.auth_user_id;

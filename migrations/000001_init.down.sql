-- Rollback initial database schema for 4Planet backend
-- Generated from GORM models to ensure perfect alignment

-- Drop views
DROP VIEW IF EXISTS user_stats;

-- Drop indexes
DROP INDEX IF EXISTS idx_webhook_events_provider;
DROP INDEX IF EXISTS idx_share_tokens_user;
DROP INDEX IF EXISTS idx_donations_project;
DROP INDEX IF EXISTS idx_donations_user_created;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_user;
DROP INDEX IF EXISTS idx_subscriptions_user;
DROP INDEX IF EXISTS idx_news_published;
DROP INDEX IF EXISTS idx_media_project;
DROP INDEX IF EXISTS idx_projects_geo;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_password_reset_user;
DROP INDEX IF EXISTS idx_email_verif_user;
DROP INDEX IF EXISTS idx_sessions_valid;
DROP INDEX IF EXISTS idx_sessions_user;

-- Drop tables
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS share_tokens;
DROP TABLE IF EXISTS donations;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS achievements;
DROP TABLE IF EXISTS news;
DROP TABLE IF EXISTS media_files;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS tree_prices;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS media_kind;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS payment_provider;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS news_type;
DROP TYPE IF EXISTS project_status;

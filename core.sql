-- Extensions (optional)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =============================
-- ENUM TYPES
-- =============================

CREATE TYPE project_status      AS ENUM ('planned','in_progress','completed');
CREATE TYPE news_type           AS ENUM ('achievement','invite','update');
CREATE TYPE payment_status      AS ENUM ('pending','succeeded','failed','refunded','canceled');
CREATE TYPE payment_provider    AS ENUM ('cloudpayments','kaspi','paypal','tribute');
CREATE TYPE subscription_status AS ENUM ('active','past_due','canceled','paused','incomplete');
CREATE TYPE user_status         AS ENUM ('pending','active','blocked');
CREATE TYPE media_kind          AS ENUM ('image','video','document');

-- =============================
-- CORE TABLES
-- =============================

-- Пользователи (внутренний менеджмент регистрации/логина)
CREATE TABLE users (
  auth_user_id     text PRIMARY KEY,
  username         text UNIQUE,                    -- публичный ник; можно NULL до установки
  display_name     text,                           -- реальное имя (можно не показывать публично)
  avatar_url       text,
  email            text UNIQUE,                    -- e-mail для логина/уведомлений
  email_verified_at timestamptz,                   -- отметка подтверждения e-mail
  status           user_status NOT NULL DEFAULT 'pending',
  password_hash    text,                           -- для локной аутентификации (bcrypt/argon2); NULL если соцлогин
  total_trees      int  NOT NULL DEFAULT 0,
  donations_count  int  NOT NULL DEFAULT 0,
  last_donation_at timestamptz,
  created_at       timestamptz NOT NULL DEFAULT now()
);

-- Простейшие сессии (cookie с opaque session_id)
CREATE TABLE sessions (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_user_id text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  created_at   timestamptz NOT NULL DEFAULT now(),
  expires_at   timestamptz NOT NULL,
  revoked_at   timestamptz,
  user_agent   text,
  ip_addr      inet
);
CREATE INDEX sessions_user_idx ON sessions(auth_user_id);
CREATE INDEX sessions_valid_idx ON sessions(expires_at) WHERE revoked_at IS NULL;

-- Токены подтверждения e-mail
CREATE TABLE email_verification_tokens (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_user_id text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  token        text UNIQUE NOT NULL,              -- случайный одноразовый токен
  created_at   timestamptz NOT NULL DEFAULT now(),
  expires_at   timestamptz NOT NULL,
  used_at      timestamptz
);
CREATE INDEX email_verif_user_idx ON email_verification_tokens(auth_user_id);

-- Токены сброса пароля
CREATE TABLE password_reset_tokens (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_user_id text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  token        text UNIQUE NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  expires_at   timestamptz NOT NULL,
  used_at      timestamptz
);
CREATE INDEX password_reset_user_idx ON password_reset_tokens(auth_user_id);

-- Фиксированная цена за одно дерево по валютам
CREATE TABLE tree_prices (
  currency    text PRIMARY KEY,      -- ISO-4217, e.g. RUB, KZT, USD
  price_minor bigint NOT NULL,       -- минимальные единицы (копейки/тиын/центы)
  updated_at  timestamptz NOT NULL DEFAULT now()
);

-- Проекты (бывш. planting_events) — для посадок, уборок рек и пр.
CREATE TABLE projects (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title            text NOT NULL,
  description      text,
  status           project_status NOT NULL DEFAULT 'planned',
  starts_at        timestamptz,
  ends_at          timestamptz,
  country_code     text,
  region           text,
  location_geojson jsonb NOT NULL,   -- GeoJSON Point/Polygon/MultiPolygon
  trees_target     integer,          -- релевантно для лесопосадок; можно NULL для других типов
  trees_planted    integer,
  created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX projects_status_idx ON projects(status);
CREATE INDEX projects_geo_idx    ON projects USING GIN (location_geojson);

-- Медиафайлы, привязанные к проектам
CREATE TABLE media_files (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id   uuid REFERENCES projects(id) ON DELETE CASCADE,
  kind         media_kind NOT NULL DEFAULT 'image',
  url          text NOT NULL,         -- абсолютная ссылка (S3/Cloud, CDN и т.п.)
  mime_type    text,
  title        text,
  alt_text     text,
  meta         jsonb DEFAULT '{}'::jsonb,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX media_project_idx ON media_files(project_id);

-- Новости/лента
CREATE TABLE news (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  type         news_type NOT NULL,
  title        text NOT NULL,
  body_md      text,
  cover_url    text,
  project_id   uuid REFERENCES projects(id) ON DELETE SET NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  published_at timestamptz
);
CREATE INDEX news_published_idx ON news(published_at DESC);

-- Ачивки (обобщение бейджей; порог по деревьям может быть NULL для кастомных выдач)
CREATE TABLE achievements (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code            text UNIQUE NOT NULL,
  title           text NOT NULL,
  description     text,
  threshold_trees integer,          -- NULL => вручную присваиваемая/кастомная ачивка
  image_url       text
);

CREATE TABLE user_achievements (
  auth_user_id text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  achievement_id uuid NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
  awarded_at   timestamptz NOT NULL DEFAULT now(),
  reason       text,                 -- опционально: причина/комментарий (например, волонтерство)
  PRIMARY KEY (auth_user_id, achievement_id)
);

-- Подписки у провайдеров (мы их отображаем/сопоставляем)
CREATE TABLE subscriptions (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_user_id             text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  provider                 payment_provider NOT NULL,
  provider_customer_id     text,
  provider_subscription_id text UNIQUE,
  amount_minor             bigint NOT NULL,
  currency                 text   NOT NULL,
  interval_months          integer NOT NULL DEFAULT 1,
  status                   subscription_status NOT NULL,
  started_at               timestamptz NOT NULL DEFAULT now(),
  canceled_at              timestamptz,
  meta                     jsonb DEFAULT '{}'::jsonb
);
CREATE INDEX subscriptions_user_idx ON subscriptions(auth_user_id);

-- Платёжные события (нормализованная платёжка отдельно от бизнес-логики)
CREATE TABLE payments (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  provider             payment_provider NOT NULL,
  provider_payment_id  text UNIQUE,             -- внешний ID транзакции
  auth_user_id         text REFERENCES users(auth_user_id) ON DELETE SET NULL, -- nullable: может быть «осиротевшим»
  subscription_id      uuid REFERENCES subscriptions(id) ON DELETE SET NULL,  -- если это чардж подписки
  amount_minor         bigint NOT NULL,
  currency             text   NOT NULL,
  status               payment_status NOT NULL,
  occurred_at          timestamptz,              -- время у провайдера
  meta                 jsonb DEFAULT '{}'::jsonb,
  created_at           timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX payments_user_idx   ON payments(auth_user_id, created_at DESC);
CREATE INDEX payments_status_idx ON payments(status);

-- Бизнес-запись «донат» — создаётся ТОЛЬКО на успешную оплату
-- trees_count вычислен по таблице tree_prices на момент создания записи
CREATE TABLE donations (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_user_id  text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  payment_id    uuid NOT NULL UNIQUE REFERENCES payments(id) ON DELETE RESTRICT,
  project_id    uuid REFERENCES projects(id) ON DELETE SET NULL,   -- nullable: донат в конкретный проект или общий
  trees_count   integer NOT NULL,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX donations_user_created_idx ON donations(auth_user_id, created_at DESC);
CREATE INDEX donations_project_idx      ON donations(project_id);

-- Шаринг короткими ссылками (OG-превью и удобные ссылки в соцсети)
CREATE TABLE share_tokens (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_user_id text NOT NULL REFERENCES users(auth_user_id) ON DELETE CASCADE,
  kind         text NOT NULL,            -- 'profile' | 'donation'
  ref_id       uuid,                     -- donations.id если kind='donation'
  slug         text UNIQUE NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX share_tokens_user_idx ON share_tokens(auth_user_id);

-- Вебхуки провайдеров (идемпотентность и аудит)
CREATE TABLE webhook_events (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  provider          payment_provider NOT NULL,
  event_type        text NOT NULL,
  event_idempotency text UNIQUE,
  received_at       timestamptz NOT NULL DEFAULT now(),
  raw_payload       jsonb NOT NULL,
  signature_ok      boolean NOT NULL,
  processed_ok      boolean NOT NULL DEFAULT false,
  processing_error  text
);
CREATE INDEX webhook_events_provider_idx ON webhook_events(provider, received_at DESC);

-- =============================
-- VIEWS
-- =============================

-- Сводная статистика по пользователю (для проверки консистентности)
CREATE OR REPLACE VIEW user_stats AS
SELECT
  u.auth_user_id,
  COALESCE(SUM(d.trees_count),0)::int AS total_trees,
  COUNT(d.*)::int                     AS donations_count,
  MAX(d.created_at)                   AS last_donation_at
FROM users u
LEFT JOIN donations d ON d.auth_user_id = u.auth_user_id
GROUP BY u.auth_user_id;

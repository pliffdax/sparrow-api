CREATE TABLE IF NOT EXISTS users (
  id          BIGSERIAL PRIMARY KEY,
  name        TEXT        NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
  id          BIGSERIAL PRIMARY KEY,
  title       TEXT        NOT NULL,
  is_global   BOOLEAN     NOT NULL DEFAULT FALSE,
  owner_id    BIGINT      NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_categories_global_title
  ON categories (title) WHERE is_global = TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS uq_categories_owner_title
  ON categories (owner_id, title) WHERE is_global = FALSE;

CREATE TABLE IF NOT EXISTS records (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  category_id  BIGINT        NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
  amount       NUMERIC(14,2) NOT NULL CHECK (amount > 0),
  happened_at  TIMESTAMPTZ   NOT NULL,
  note         TEXT NULL,
  created_at   TIMESTAMPTZ   NOT NULL DEFAULT now()
);

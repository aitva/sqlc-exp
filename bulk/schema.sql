-- Alias introduced to workaround sqlc bug (#1851)
CREATE DOMAIN public.TEXT_NULL AS TEXT;

-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT      NOT NULL,
  bio  TEXT
);

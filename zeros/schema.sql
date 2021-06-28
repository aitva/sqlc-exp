-- Example queries for sqlc
CREATE TABLE authors (
  id    BIGSERIAL PRIMARY KEY,
  name  text      NOT NULL,
  birth date      NOT NULL,
  bio   text
);

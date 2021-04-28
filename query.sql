-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: CreateAuthors :exec
INSERT INTO authors
SELECT unnest(@ids::bigint[]) AS id,
  unnest(@names::text[]) AS name,
  unnest(@bios::text[]) AS bio;

-- name: UpdateAuthors :exec
UPDATE authors AS a
SET name = tmp.name, bio = tmp.bio
FROM (
  SELECT unnest(@ids::bigint[]) AS id,
  unnest(@names::text[]) AS name,
  unnest(@bios::text[]) AS bio
) AS tmp
WHERE s.id = tmp.id;

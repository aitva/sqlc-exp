-- name: CreateAuthor :exec
INSERT INTO authors (id, name, bio) VALUES ($1, $2, $3);

-- name: CreateAuthors :exec
INSERT INTO authors
SELECT unnest(@ids::BIGINT[]) AS id,
  unnest(@names::TEXT[]) AS name,
  unnest(@bios::TEXT_NULL[]) AS bio;

-- name: UpdateAuthor :exec
UPDATE authors SET name = $2, bio = $3 WHERE id = $1;

-- name: UpdateAuthors :exec
UPDATE authors AS a
SET name = tmp.name, bio = tmp.bio
FROM (
  SELECT unnest(@ids::BIGINT[]) AS id,
  unnest(@names::TEXT[]) AS name,
  unnest(@bios::TEXT_NULL[]) AS bio
) AS tmp
WHERE a.id = tmp.id;

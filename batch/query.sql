-- name: CreateAuthor :exec
INSERT INTO authors (id, name, bio) VALUES ($1, $2, $3);

-- name: CreateAuthors :exec
INSERT INTO authors
SELECT unnest(@ids::bigint[]) AS id,
  unnest(@names::text[]) AS name,
  unnest(@bios::text[]) AS bio;

-- name: UpdateAuthor :exec
UPDATE authors SET name = $2, bio = $3 WHERE id = $1;

-- name: UpdateAuthors :exec
UPDATE authors AS a
SET name = tmp.name, bio = tmp.bio
FROM (
  SELECT unnest(@ids::bigint[]) AS id,
  unnest(@names::text[]) AS name,
  unnest(@bios::text[]) AS bio
) AS tmp
WHERE s.id = tmp.id;

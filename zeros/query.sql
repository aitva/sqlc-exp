-- name: CreateAuthor :one
INSERT INTO authors (id, name, birth, bio) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListAuthors :many
SELECT *
FROM authors
WHERE (@name::text = '' OR name = @name)
    AND (@birth::date = '0001-01-01' OR birth = @birth)
    AND (@bio::text = '' OR bio ILIKE '%' || @bio || '%');

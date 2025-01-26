-- name: CreateBook :one
INSERT INTO Book (title, author_id) VALUES ($1, $2) RETURNING *;

-- name: ReadBook :one
SELECT * FROM Book WHERE id = $1;

-- name: ReadBooks :many
SELECT * FROM Book;

-- name: UpdateBook :exec
UPDATE Book SET title = $1, author_id = $2 WHERE id = $3;

-- name: DeleteBook :exec
DELETE FROM Book WHERE id = $1;

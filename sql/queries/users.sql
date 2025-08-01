-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
select * from users
where name = $1
limit 1;

-- name: Resets :exec
delete from users;

-- name: GetUsers :many
select name from users;
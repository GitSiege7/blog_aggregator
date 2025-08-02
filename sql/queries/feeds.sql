-- name: CreateFeed :exec
insert into feeds (id, created_at, updated_at, name, url, user_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: GetFeeds :many
select name, url, user_id from feeds;
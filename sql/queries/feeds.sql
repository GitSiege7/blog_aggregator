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

-- name: CreateFeedFollow :one
with new_feed_follow as (
    insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4, $5)
    returning *
) select new_feed_follow.*, feeds.name as feed_name, users.name as user_name
from new_feed_follow
join feeds on new_feed_follow.feed_id = feeds.id
join users on new_feed_follow.user_id = users.id;

-- name: GetFeedByUrl :one
select * from feeds
where feeds.url = $1;

-- name: GetFeedFollowsForUser :many
select
    feed_follows.*,
    feeds.name as feed_name,
    users.name as user_name
from feed_follows
join feeds on feed_follows.feed_id = feeds.id
join users on feed_follows.user_id = users.id
where feed_follows.user_id = $1;
-- name: CreatePost :exec
insert into posts (id, created_at, updated_at, title, url, description, published_at, feed_id) 
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: GetPostsForUser :many
select posts.*, feeds.name as feed_name from posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
order by posts.published_at desc
limit $2;
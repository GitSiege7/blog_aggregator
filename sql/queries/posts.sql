-- name: CreatePost :exec
insert into posts values (
    $1,
    $2,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)

-- name: GetPostsForUser :many
select * from posts
order by posts.created_at asc
limit $1;
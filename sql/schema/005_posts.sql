-- +goose Up
create table posts (
    id uuid primary key not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text,
    url text not null,
    description text,
    published_at timestamp not null,
    feed_id uuid not null,
    foreign key (feed_id) references feeds.id,
    unique(url)
);

-- +goose Down
drop table posts;
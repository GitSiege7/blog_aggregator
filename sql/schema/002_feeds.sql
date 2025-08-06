-- +goose Up
create table feeds(
    id UUID primary key not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text not null,
    user_id UUID not null references users(id) on delete cascade,
    unique(url)
);

-- +goose Down
drop table feeds;
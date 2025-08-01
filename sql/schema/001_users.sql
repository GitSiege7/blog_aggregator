-- +goose Up
create table users(
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null
);

-- +goose Down
drop table users;
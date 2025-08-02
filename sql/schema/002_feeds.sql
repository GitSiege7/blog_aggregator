-- +goose Up
create table feeds(
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text not null,
    user_id UUID not null,
    unique(url),
    foreign key(user_id) references users(id) on delete cascade
);

-- +goose Down
drop table feeds;
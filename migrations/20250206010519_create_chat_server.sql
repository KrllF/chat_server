-- +goose Up
create table chats (
    id serial primary key,
    members text not null,
    created_at timestamp not null default now()
);

CREATE TABLE messages (
    sender text not null,
    letter text not null,
    delivery_time timestamp not null default now()
);

-- +goose Down
drop table chats; -- Удаляем таблицу chats
drop table messages; -- Удаляем таблицу messages
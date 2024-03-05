-- +goose Up
create table chats(
    id serial primary key
);
create table chats_users(
  id serial primary key,
  chat_id   integer not null,
  user_name text not null,
  foreign key (chat_id) references chats(id) on delete cascade
);
create table messages(
    id serial primary key,
    chat_id integer not null,
    from_user text not null,
    text text not null,
    timestamp timestamp not null default now(),
    foreign key (chat_id) references chats(id) on delete cascade
);

-- +goose Down
drop table chats_users;
drop table messages;
drop table chats;
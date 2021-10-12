-- migrate:up
create table "user"(
  id serial primary key,
  email text not null unique,
  password_hash bytea not null,
  name text not null default ''
);

-- migrate:down
drop table "user";

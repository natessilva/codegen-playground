-- migrate:up
create table identity(
  id serial primary key,
  email text not null unique,
  password_hash bytea not null,
  name text not null default ''
);

create table space(
  id serial primary key,
  name text not null default ''
);

create table "user"(
  space_id int not null references space (id) on delete cascade,
  identity_id int not null references identity(id) on delete cascade,
  primary key(space_id, identity_id)
);

alter table identity add column current_space_id int references space(id) on delete set null;

-- migrate:down
drop table "user";
drop table space;
drop table identity;

-- name: CreateUser :one
insert into "user"(email, password_hash)
values($1,$2)
on conflict do nothing
returning id;

-- name: GetUser :one
select
  name
from "user"
where id = $1;

-- name: GetPasswordById :one
select
  password_hash
from "user"
where id = $1;

-- name: GetPasswordByEmail :one
select
  id,
  password_hash
from "user"
where email = $1;

-- name: SetPassword :exec
update "user"
set password_hash = $1
where id = $2;

-- name: UpdateUser :exec
update "user"
set name = $2
where id = $1;
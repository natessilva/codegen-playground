-- name: Create :one
insert into "user"(email, password_hash)
values($1,$2)
on conflict do nothing
returning id;

-- name: Get :one
select
  *
from "user"
where id = $1;

-- name: GetByEmail :one
select
  *
from "user"
where email = $1;

-- name: Update :exec
update "user"
set name = $2
where id = $1;

-- name: SetPassword :exec
update "user"
set password_hash = $2
where id = $1;
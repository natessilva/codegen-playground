-- name: CreateSpace :one
insert into space(name)
values($1)
returning id;

-- name: GetSpace :one
select
  *
from space
where id = $1;

-- name: UpdateSpace :exec
update space
set name = $2
where id = $1;
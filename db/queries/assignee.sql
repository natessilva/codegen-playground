-- name: Assign :exec
insert into assignee(space_id, ticket_id, identity_id)
values($1, $2, $3)
on conflict do nothing;
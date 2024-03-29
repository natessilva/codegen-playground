// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: space.sql

package model

import (
	"context"
)

const createSpace = `-- name: CreateSpace :one
insert into space(name)
values($1)
returning id
`

func (q *Queries) CreateSpace(ctx context.Context, name string) (int32, error) {
	row := q.db.QueryRow(ctx, createSpace, name)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getSpace = `-- name: GetSpace :one
select
  id, name
from space
where id = $1
`

func (q *Queries) GetSpace(ctx context.Context, id int32) (Space, error) {
	row := q.db.QueryRow(ctx, getSpace, id)
	var i Space
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listSpaces = `-- name: ListSpaces :many
select s.id, s.name
from "user" u
join space s on u.space_id = s.id
where u.identity_id = $1
`

func (q *Queries) ListSpaces(ctx context.Context, identityID int32) ([]Space, error) {
	rows, err := q.db.Query(ctx, listSpaces, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Space
	for rows.Next() {
		var i Space
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSpace = `-- name: UpdateSpace :exec
update space
set name = $2
where id = $1
`

type UpdateSpaceParams struct {
	ID   int32
	Name string
}

func (q *Queries) UpdateSpace(ctx context.Context, arg UpdateSpaceParams) error {
	_, err := q.db.Exec(ctx, updateSpace, arg.ID, arg.Name)
	return err
}

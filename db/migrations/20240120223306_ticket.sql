-- migrate:up
create table ticket(
    space_id int references space(id) on delete cascade,
    id serial,
    primary key(space_id, id),

    subject text not null,
    body text not null default ''
);

-- migrate:down
drop table ticket;

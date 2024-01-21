-- migrate:up
create table assignee(
    space_id int references space(id) on delete cascade,
    ticket_id int,
    identity_id int,

    primary key(space_id, ticket_id, identity_id),
    foreign key(space_id, ticket_id) references ticket (space_id, id) on delete cascade,
    foreign key(space_id, identity_id) references "user"(space_id, identity_id) on delete cascade
);

create index on assignee(identity_id);

-- migrate:down
drop table assignee;
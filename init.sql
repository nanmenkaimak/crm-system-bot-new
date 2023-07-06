create table if not exists users
(
    id         serial primary key,
    name       varchar(50) not null,
    username   varchar(50) not null,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists rooms
(
    id         serial primary key,
    admin_id   int references users (id) on delete cascade not null,
    name       varchar(50),
    capacity   int,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists beds
(
    id         serial primary key,
    room_id    int references rooms (id) on delete cascade not null,
    name       varchar(50),
    cost       int,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

create table if not exists customers
(
    id            serial,
    bed_id        int references beds (id) on delete cascade not null,
    full_name     varchar(100),
    photo         varchar(200),
    phone         varchar(50),
    info          varchar(200),
    money         int,
    is_here       boolean default true,
    arrival_day   timestamp without time zone,
    departure_day timestamp without time zone,
    created_at    timestamp without time zone,
    updated_at    timestamp without time zone
);

create table if not exists expenses
(
    id         serial,
    name       varchar(100),
    money      int,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


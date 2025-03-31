-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id          uuid      default gen_random_uuid() not null
        constraint users_pk
            primary key,
    first_name  varchar(100)                        not null,
    second_name varchar(100)                        not null,
    birthdate   date,
    city        varchar(100),
    biography   text,
    password    varchar(60)                         not null,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    updated_at  timestamp default CURRENT_TIMESTAMP not null
);

-- +goose Down
DROP TABLE users;
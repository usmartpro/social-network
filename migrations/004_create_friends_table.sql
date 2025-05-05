-- +goose Up
CREATE TABLE IF NOT EXISTS friends
(
    id         uuid      default gen_random_uuid() not null
        constraint friends_pk
            primary key,
    user_id1   uuid                                not null
        constraint friends_users_id1_fk
            references public.users
            on update cascade on delete cascade,
    user_id2   uuid                                not null
        constraint friends_users_id2_fk
            references public.users
            on update cascade on delete cascade,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    constraint table_name_pk
        unique (user_id1, user_id2)
);

-- +goose Down
DROP TABLE friends;
-- +goose Up
CREATE TABLE IF NOT EXISTS posts
(
    id             uuid      default gen_random_uuid() not null
        constraint posts_pk
            primary key,
    text           text                                not null,
    author_user_id uuid                                not null
        constraint posts_users_id_fk
            references public.users
            on update cascade on delete cascade,
    created_at     timestamp default CURRENT_TIMESTAMP not null,
    updated_at     timestamp default CURRENT_TIMESTAMP not null
);

-- +goose Down
DROP TABLE posts;
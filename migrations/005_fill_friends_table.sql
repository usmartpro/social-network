-- +goose Up
INSERT INTO public.friends (user_id1, user_id2)
WITH first_10_users AS (
    SELECT id
    FROM public.users
    ORDER BY id
    LIMIT 10
)
SELECT
    u.id AS user_id,
    f.friend_id
FROM
    first_10_users u
        CROSS JOIN LATERAL (
        SELECT friend_id
        FROM (
                 SELECT id AS friend_id
                 FROM public.users
                 WHERE id != u.id
             ) all_friends
        ORDER BY RANDOM()
        LIMIT 1000
        ) f;

-- +goose Down

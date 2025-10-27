-- Creates table to track users

-- +goose Up

BEGIN;

CREATE TABLE IF NOT EXISTS app.users
(
    id bigint DEFAULT random(1, 9223372036854775807),
    email character varying(320),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    UNIQUE(email)
);

COMMIT;

-- +goose Down

BEGIN;

DROP TABLE IF EXISTS app.users;

COMMIT;


-- Creates table to track users

-- +goose Up

BEGIN;

CREATE TABLE IF NOT EXISTS app.users
(
    id bigint DEFAULT random(1, 9223372036854775807),
    google_sub character varying(40),
    email character varying(320),
    picture_url bpchar,
    full_name character varying(255),
    given_name character varying(255),
    family_name character varying(255),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    UNIQUE(email),
    UNIQUE(sub)
);

COMMIT;

-- +goose Down

BEGIN;

DROP TABLE IF EXISTS app.users;

COMMIT;


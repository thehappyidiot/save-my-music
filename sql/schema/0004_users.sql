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
    spotify_id character varying(32),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    FOREIGN KEY (spotify_id)
        REFERENCES spotify.users (spotify_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    UNIQUE(email),
    UNIQUE(google_sub)
);

COMMIT;

-- +goose Down

BEGIN;

DROP TABLE IF EXISTS app.users;

COMMIT;


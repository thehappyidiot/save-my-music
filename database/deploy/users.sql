-- Deploy save-my-music:users to pg
-- requires: appschema

BEGIN;

CREATE TABLE IF NOT EXISTS app.users
(
    id bigint DEFAULT random(1, 9223372036854775807),
    email character varying(320),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

COMMIT;

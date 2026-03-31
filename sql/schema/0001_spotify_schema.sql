-- Schema for data from Spotify

-- +goose Up

BEGIN;

CREATE SCHEMA spotify;

COMMIT;

-- +goose Down

BEGIN;

DROP SCHEMA IF EXISTS spotify;

COMMIT;
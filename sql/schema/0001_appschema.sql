-- Schema for all common objects, not belonging to a specific music provider

-- +goose Up

BEGIN;

CREATE SCHEMA app; 

COMMIT;

-- +goose Down

BEGIN;

DROP SCHEMA IF EXISTS app;

COMMIT;

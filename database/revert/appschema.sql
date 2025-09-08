-- Revert save-my-music:appschema from pg

BEGIN;

DROP SCHEMA IF EXISTS app;

COMMIT;

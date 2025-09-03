-- Revert save-my-music:appschema from pg

BEGIN;

DROP SCHEMA app;

COMMIT;

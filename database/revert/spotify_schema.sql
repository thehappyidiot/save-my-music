-- Revert save-my-music:spotify_schema from pg

BEGIN;

DROP SCHEMA IF EXISTS spotify;

COMMIT;

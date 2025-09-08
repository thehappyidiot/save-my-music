-- Revert save-my-music:users from pg

BEGIN;

DROP TABLE IF EXISTS app.users;

COMMIT;

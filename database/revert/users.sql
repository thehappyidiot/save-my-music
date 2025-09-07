-- Revert save-my-music:users from pg

BEGIN;

SELECT id, email 
    FROM app.users
    WHERE FALSE;

COMMIT;

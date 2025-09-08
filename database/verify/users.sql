-- Verify save-my-music:users on pg

BEGIN;

SELECT id, email 
    FROM app.users
    WHERE FALSE;

ROLLBACK;

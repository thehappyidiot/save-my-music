-- Verify save-my-music:appschema on pg

BEGIN;

SELECT pg_catalog.has_schema_privilege('app', 'usage');

ROLLBACK;
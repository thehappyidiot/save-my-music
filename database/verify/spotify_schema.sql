-- Verify save-my-music:spotify_schema on pg

BEGIN;

SELECT pg_catalog.has_schema_privilege('spotify', 'usage');

ROLLBACK;
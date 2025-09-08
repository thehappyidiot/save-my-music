-- Revert save-my-music:spotify_tables from pg

BEGIN;

DROP TABLE IF EXISTS spotify.playlist_items;
DROP TYPE IF EXISTS spotify.item_type;
DROP TABLE IF EXISTS spotify.episodes;
DROP TYPE IF EXISTS spotify.show;
DROP TABLE IF EXISTS spotify.tracks;
DROP TYPE IF EXISTS spotify.artist;
DROP TABLE IF EXISTS spotify.playlists;
DROP TABLE IF EXISTS spotify.users;

COMMIT;

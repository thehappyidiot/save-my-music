-- Verify save-my-music:spotify_tables on pg

BEGIN;

SELECT 
    spotify_id, display_name, internal_id, last_scan, created
    FROM spotify.users
    WHERE false;

SELECT 
    owner, playlist_id, name, description, created, is_deleted
    FROM spotify.playlists
    WHERE false;

SELECT 
    track_id, album_id, album_name, (artists[0]).name, (artists[0]).artist_id, 
    available_markets, is_explicit, isrc, upc, restriction_reason,
    name, is_local
    FROM  spotify.tracks
    WHERE false;

SELECT
    episode_id, is_explicit, description, duration_ms, languages, name,
    release_date, restriction_reason, (show).show_id, (show).available_markets,
    (show).description, (show).is_explicit, (show).languages, (show).name, 
    (show).publisher
    FROM spotify.episodes;

SELECT 
    playlist_id, track_id, item_type, added_at, added_by, is_local,
    is_unavailable, is_deleted, episode_id
    FROM spotify.playlist_items;

ROLLBACK;

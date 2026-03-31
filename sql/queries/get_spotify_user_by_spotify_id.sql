-- name: GetSpotifyUserBySpotifyId :one

SELECT
    spotify_id,
    display_name,
    last_scan,
    created_at
FROM
    spotify.users su
WHERE
    su.spotify_id = $1; 
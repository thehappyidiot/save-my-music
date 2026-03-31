-- name: UpsertSpotifyUser :one

INSERT INTO spotify.users (
        spotify_id,
        display_name,
        created_at,
        updated_at
    )
VALUES (
    $1,
    $2,
    NOW(),
    NOW()
) ON CONFLICT (spotify_id) DO
UPDATE
SET display_name = $2,
    updated_at = NOW()
RETURNING *; 
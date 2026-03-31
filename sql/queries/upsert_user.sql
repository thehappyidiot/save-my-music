-- name: UpsertUser :one
INSERT INTO app.users (
        google_sub,
        email,
        picture_url,
        full_name,
        given_name,
        family_name,
        created_at,
        updated_at,
        spotify_id
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        NOW(),
        NOW(),
        $7
    ) ON CONFLICT (google_sub) DO
UPDATE
SET email = $2,
    picture_url = $3,
    full_name = $4,
    given_name = $5,
    family_name = $6,
    updated_at = NOW(),
    spotify_id = $7
RETURNING *;
-- name: GetUserBySub :one

SELECT 
    id,
    google_sub,
    email,
    picture_url,
    full_name,
    given_name,
    family_name,
    created_at,
    updated_at
FROM 
    app.users u
WHERE
    u.google_sub = $1;
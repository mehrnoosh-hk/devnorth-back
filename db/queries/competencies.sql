-- name: CreateCompetency :one
INSERT INTO competencies (
    name,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetCompetencyByName :one
SELECT * FROM competencies
WHERE name = $1
LIMIT 1;

-- name: GetCompetencyByID :one
SELECT * FROM competencies
WHERE id = $1
LIMIT 1;

-- name: UpdateCompetencyDescription :one
UPDATE competencies
SET description = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: GetAllCompetencies :many
SELECT * FROM competencies
ORDER BY created_at DESC;

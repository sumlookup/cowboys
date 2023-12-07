-- name: CreateCowboy :one
INSERT INTO cowboys (name, health, damage) VALUES (@name, @health, @damage) RETURNING id::uuid;

-- name: CreateManyCowboys :batchmany
INSERT INTO cowboys (name, health, damage) VALUES (@name, @health, @damage) RETURNING *;

-- name: UpdateCowboyHealth :one
UPDATE cowboys SET health = (health + @health),
    updated_at = NOW() WHERE id = @id::uuid RETURNING *;

-- name: DeleteAllCowboys :one
DELETE from cowboys RETURNING *;

-- name: GetSingleCowboyByName :one
SELECT * from cowboys where name = @name limit 1;

-- name: ListAliveCowboys :many
SELECT * FROM cowboys WHERE health > 0 AND
    deleted_at IS NULL
ORDER BY
    CASE WHEN @query_sort::text = 'desc' THEN created_at END DESC,
    CASE WHEN @query_sort::text = 'asc' THEN created_at END ASC
LIMIT @query_limit::int
    OFFSET @query_offset::int;

-- name: GetSingleAliveCowboyById :one
SELECT * from cowboys where id = @id AND health > 0;

-- name: GetRandomCowboy :one
SELECT * from cowboys
Where health > 0 AND id != @id
order by random()
LIMIT 1;

-- name: ReduceCowboyHealth :exec
UPDATE cowboys SET health = (health - @shooter_damage),
updated_at = NOW() WHERE id = @victim_id::uuid;

-- name: GetSingleAliveCowboyAndCount :one
SELECT a.id, a.created_at, a.name, a.health, a.damage,
       (
           Select count(*) from cowboys as b where health > 0 AND deleted_at IS NULL AND b.id != @id
       ) AS available
from cowboys as a where a.id = @id AND health > 0;

-- name: GetSingleCowboyHealth :one
SELECT health from cowboys where id = @id;
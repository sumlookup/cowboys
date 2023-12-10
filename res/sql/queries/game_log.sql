-- name: CreateGameLog :exec
INSERT INTO game_logs (game_id, shooter_id, receiver_id, damage, receiver_health, shooter_health) VALUES
    (@game_id, @shooter_id, @receiver_id, @damage, @receiver_health, @shooter_health);

-- name: ListGameLogsByGameId :many
SELECT * FROM game_logs WHERE game_id = @game_id
ORDER BY
    CASE WHEN @query_sort::text = 'desc' THEN created_at END DESC,
    CASE WHEN @query_sort::text = 'asc' THEN created_at END ASC
LIMIT @query_limit::int
OFFSET @query_offset::int;

-- name: CountAllGameLogs :one
SELECT count(*) as total_count FROM game_logs WHERE game_id = @game_id;

SELECT a.id, a.created_at, a.name, a.health, a.damage,
       (
           Select count(*) from cowboys as b where health > 0 AND deleted_at IS NULL AND b.id != @id
       ) AS available
from cowboys as a where a.id = @id AND health > 0;
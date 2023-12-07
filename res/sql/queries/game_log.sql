-- name: CreateGameLog :one
INSERT INTO game_logs (game_id, shooter, receiver, damage) VALUES (@game_id, @shooter, @receiver, @damage) RETURNING id::uuid;

-- name: ListGameLogsByGameId :many
SELECT * FROM game_logs WHERE game_id = @game_id AND
    deleted_at IS NULL
ORDER BY
    CASE WHEN @query_sort::text = 'desc' THEN created_at END DESC,
    CASE WHEN @query_sort::text = 'asc' THEN created_at END ASC
LIMIT @query_limit::int
    OFFSET @query_offset::int;



-- name: CreateGame :one
INSERT INTO games (mode) VALUES (@mode) RETURNING id::uuid;

-- name: GetGameById :one
SELECT * FROM games WHERE id = @id AND deleted_at IS NULL;

-- name: UpdateGameWinner :exec
UPDATE games SET winner = @winner, winner_id = @winner_id::uuid, ended = true,
    updated_at = NOW() WHERE id = @game_id::uuid;


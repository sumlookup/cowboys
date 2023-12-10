-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS cowboys
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    name TEXT NOT NULL UNIQUE,
    health INT NOT NULL,
    damage INT NOT NULL
);

CREATE TABLE IF NOT EXISTS games
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    mode TEXT NOT NULL default '',
    winner TEXT NOT NULL default '',
    winner_id UUID DEFAULT NULL,
    ended bool NOT NULL default false
);


CREATE TABLE IF NOT EXISTS game_logs
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    game_id UUID NOT NULL,
    shooter_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    damage INT NOT NULL,
    receiver_health INT NOT NULL,
    shooter_health INT NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE cowboys;
DROP TABLE games;
DROP TABLE game_logs;
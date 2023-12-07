// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Cowboy struct {
	ID        uuid.UUID          `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	Name      string             `json:"name"`
	Health    int32              `json:"health"`
	Damage    int32              `json:"damage"`
}

type Game struct {
	ID        uuid.UUID          `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	Mode      string             `json:"mode"`
	Winner    string             `json:"winner"`
	WinnerID  pgtype.UUID        `json:"winner_id"`
	Ended     bool               `json:"ended"`
}

type GameLog struct {
	ID        uuid.UUID          `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	GameID    uuid.UUID          `json:"game_id"`
	Shooter   uuid.UUID          `json:"shooter"`
	Receiver  uuid.UUID          `json:"receiver"`
	Damage    int32              `json:"damage"`
}
package engine

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	daodb "github.com/sumlookup/cowboys/dao/db"
	basic "github.com/sumlookup/cowboys/engine/basic"
	intermediate "github.com/sumlookup/cowboys/engine/intermediate"
	"github.com/sumlookup/mini/service"
)

var (
	ENGINE_MODE_BASIC        = "basic"
	ENGINE_MODE_INTERMEDIATE = "intermediate"
	ENGINE_MODE_ADVANCED     = "advanced" // Had big dreams with this project
)

type CowboysEngine interface {
	Run(ctx context.Context) error
	StartGame(ctx context.Context, id, gameId uuid.UUID, name string) error
	ShootRandomCowboy(ctx context.Context, shooterId, gameId uuid.UUID, shooterName string, shooterDmg int32) (int32, error)
	SetWinner(ctx context.Context, name string, winnerId, gameId uuid.UUID) error
	GameMode() string
}

func NewEngine(name string, db daodb.Querier, pool *pgxpool.Pool, service *service.Service) CowboysEngine {
	switch name {
	case ENGINE_MODE_BASIC:
		log.Info("Returning basic game engine")
		return basic.NewBasicGameEngine(db, pool)
	case ENGINE_MODE_INTERMEDIATE:
		log.Info("Returning intermediate game engine")
		return intermediate.NewIntermediateGameEngine(db, pool, service)
	default:
		log.Info("Returning default basic game engine")
		return basic.NewBasicGameEngine(db, pool)
	}
}

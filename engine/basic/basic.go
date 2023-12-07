package basic

import (
	"context"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	daodb "github.com/sumlookup/cowboys/dao/db"
	"time"
)

type BasicGameEngine struct {
	Dao  daodb.Querier
	Db   *pgxpool.Pool
	Conn *pgx.Conn
}

type Game struct {
	GameChan chan bool
	Cowboys  []*daodb.Cowboy
}

// Basic implementation :
// Data is stored and managed in Cockroachdb
// Each cowboy has its own PROCESS(goroutine) and calls an HTTP endpoint which gets routed to our exposed gateway that is pointing to grpc endpoint
// When the endpoint receives the call it uses the engine to invoke database transaction
// before every action the cowboy checks its own health and shoots random

func NewBasicGameEngine(db daodb.Querier, pool *pgxpool.Pool) *BasicGameEngine {
	return &BasicGameEngine{
		Db:  pool,
		Dao: db,
	}
}

func (b *BasicGameEngine) Run(ctx context.Context) error {

	log := logrus.New()

	mode := b.GameMode()
	log.Infof("Starting : %s mode", mode)
	dao := daodb.New(b.Db)

	// get players from DB
	cowboys, err := dao.ListAliveCowboys(ctx, daodb.ListAliveCowboysParams{
		QuerySort:   "asc", // Oldest cowboys first
		QueryOffset: 0,
		QueryLimit:  5, // 5 at a time as per requirements
	})
	if err != nil {
		log.Error("failed to to fetch alive cowboys : %v", err)
		return err
	}
	if len(cowboys) <= 1 {
		log.Error("more than one cowboy is required to simulate this shooter")
		return fmt.Errorf("more than one cowboy is required to simulate this shooter, currently found : %v", len(cowboys))
	}
	gameId, err := dao.CreateGame(ctx, b.GameMode())
	if err != nil {
		log.Error("failed to create a new game : %v", err)
		return err
	}

	log.Infof("Found %v pLayers", len(cowboys))

	for _, cowboy := range cowboys {
		go b.StartGame(context.Background(), cowboy.ID, gameId, cowboy.Name)
	}

	return nil
}

func (b *BasicGameEngine) StartGame(ctx context.Context, id, gameId uuid.UUID, name string) error {
	log := logrus.WithContext(ctx)
	log.Infof("StartGame - %v starting the game", name)
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			db := daodb.New(b.Db)
			// get current state
			cb, err := db.GetSingleAliveCowboyAndCount(ctx, id)
			if err != nil {
				if err.Error() == "no rows in result set" {
					log.Infof("%v is Dead", name)
					return nil
				}
				log.Error("failed to get alive cowboy state : ", err)
				return err
			}

			if cb.Available == 0 {
				log.Infof("%v is the Winner", name)
				b.SetWinner(ctx, name, id, gameId)
				return nil
			}
			//todo: do something with URL -- env?
			err = makeHttpRequest("http://localhost:9091/shoot", map[string]interface{}{
				"shooter_id":     id.String(),
				"shooter_name":   name,
				"shooter_damage": cb.Damage,
			})
			if err != nil {
				log.Errorf("failed while making http request - %v", err)
				return err
			}

		}
	}

}

func (b *BasicGameEngine) ShootRandomCowboy(ctx context.Context, shooterId uuid.UUID, shooterName string, shooterDmg int32) (int32, error) {
	log := logrus.WithContext(ctx)

	st := map[string]interface{}{}

	err := crdbpgx.ExecuteTx(context.Background(), b.Db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return shootAtCowboyTx(context.Background(), tx, shooterId, shooterName, shooterDmg, st)
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Warnf("%v shot final shot before dying and missed", shooterName)
			return 0, nil
		}
		log.Errorf("failed updating DB with shooter things : %v", err)
		return 0, nil
	}
	log.Infof("%v H: %v : shot at : %v, with damage : %v, and reduced health from : %v to : %v", st["shooter"], st["shooter_health"], st["victim_name"], st["shooter_damage"], st["victim_health"], st["victim_health_after"])
	val, ok := st["victim_health_after"].(int32)
	if !ok {
		log.Error("failed interface converstion : %v", err)
		return 0, nil
	}

	return val, nil
}

func (b *BasicGameEngine) SetWinner(ctx context.Context, name string, winnerId, gameId uuid.UUID) error {
	return b.Dao.UpdateGameWinner(ctx, daodb.UpdateGameWinnerParams{
		Winner:   name,
		WinnerID: winnerId,
		GameID:   gameId,
	})
}

func (b *BasicGameEngine) GameMode() string {
	return "basic"
}

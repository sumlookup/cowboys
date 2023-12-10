package basic

import (
	"context"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/sumlookup/cowboys/dao"
	daodb "github.com/sumlookup/cowboys/dao/db"
	"os"
	"time"
)

type BasicGameEngine struct {
	Dao     daodb.Querier
	Db      *pgxpool.Pool
	HttpUrl string
}

type Game struct {
	GameChan chan bool
	Cowboys  []*daodb.Cowboy
}

func NewBasicGameEngine(db daodb.Querier, pool *pgxpool.Pool) *BasicGameEngine {
	url := os.Getenv("BASIC_GAME_URL")
	if os.Getenv("BASIC_GAME_URL") == "" {
		url = "http://localhost:9090"
	}
	return &BasicGameEngine{
		HttpUrl: url,
		Db:      pool,
		Dao:     db,
	}
}

func (b *BasicGameEngine) Run(ctx context.Context) error {

	log := logrus.WithContext(ctx)

	mode := b.GameMode()
	log.Infof("Starting : %s mode", mode)

	// get players from DB
	cowboys, err := b.Dao.ListAliveCowboys(ctx, daodb.ListAliveCowboysParams{
		QuerySort:   "asc", // Oldest cowboys first
		QueryOffset: 0,
		QueryLimit:  5, // 5 at a time as per requirements
	})
	if err != nil {
		log.Errorf("failed to to fetch alive cowboys : %v", err)
		return err
	}
	if len(cowboys) <= 1 {
		log.Error("more than one cowboy is required to simulate this shooter")
		return fmt.Errorf("more than one cowboy is required to simulate this shooter, currently found : %v", len(cowboys))
	}
	gameId, err := b.Dao.CreateGame(ctx, b.GameMode())
	if err != nil {
		log.Errorf("failed to create a new game : %v", err)
		return fmt.Errorf("failed to create a new game : %v", err)
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

	for range t.C {

		// get current state
		cb, err := b.Dao.GetSingleAliveCowboyAndCount(ctx, id)
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
			err = b.SetWinner(ctx, name, id, gameId)
			if err != nil {
				log.Errorf("failed to set %v, : %v, as the winner - %v", name, id, err)
			}
			return nil
		}
		reqData := fmt.Sprintf("{\"shooter_id\":\"%s\",\"game_id\":\"%s\",\"shooter_name\":\"%s\",\"shooter_damage\":\"%v\"}",
			id.String(),
			gameId.String(),
			name,
			cb.Damage)

		err = postHttpRequest(ctx, b.HttpUrl+"/shoot", reqData)
		if err != nil {
			log.Errorf("failed while making http request - %v", err)
			return err
		}

	}

	return nil
}

func (b *BasicGameEngine) ShootRandomCowboy(ctx context.Context, shooterId, gameId uuid.UUID, shooterName string, shooterDmg int32) (int32, error) {
	log := logrus.WithContext(ctx)

	st := map[string]interface{}{}

	err := crdbpgx.ExecuteTx(ctx, b.Db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return dao.ShootAtCowboyTx(ctx, tx, shooterId, shooterName, shooterDmg, st)
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Warnf("%v shot final shot before dying and missed", shooterName)
			return 0, nil
		}
		log.Errorf("failed updating DB with shooter damage : %v", err)
		return 0, fmt.Errorf("failed updating DB with shooter damage : %v", err)
	}
	log.Infof("%v H: %v : shot at : %v, with damage : %v, and reduced health from : %v to : %v", st["shooter_name"], st["shooter_health"], st["receiver_name"], st["shooter_damage"], st["receiver_health"], st["receiver_health_after"])
	val, ok := st["receiver_health_after"].(int32)
	if !ok {
		log.Warnf("failed interface conversion : %v", err)
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

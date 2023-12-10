package basic

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/sumlookup/cowboys/dao/connection"
	"github.com/sumlookup/cowboys/mocks"
	"net/http"
	"net/http/httptest"

	//"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sumlookup/cowboys/dao/db"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
	"time"
)

var (
	shooterId uuid.UUID
	billId    uuid.UUID
	gameId    uuid.UUID
	dbPool    *pgxpool.Pool
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.ErrorLevel)
}

// setup all required data for tests
func TestMain(m *testing.M) {
	shooterId = uuid.MustParse("3bf49568-e955-4d2f-b78f-7df603571d34")
	billId = uuid.MustParse("3bf49568-e955-4d2f-b78f-7df603571d35")
	gameId = uuid.MustParse("3bf49568-e955-4d2f-b78f-7df603571d31")
	os.Exit(m.Run())
}

func Test_Run(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	cbs := []*db.Cowboy{
		{
			ID:        shooterId,
			CreatedAt: pgtype.Timestamptz{Time: time.Now().AddDate(0, 0, -1)},
			Name:      "John",
			Health:    1,
			Damage:    1,
		},
		{
			ID:        billId,
			CreatedAt: pgtype.Timestamptz{Time: time.Now()},
			Name:      "Bill",
			Health:    1,
			Damage:    1,
		},
	}

	dbmocks.EXPECT().ListAliveCowboys(gomock.Any(), db.ListAliveCowboysParams{
		QuerySort:   "asc",
		QueryOffset: 0,
		QueryLimit:  5,
	}).
		Return([]*db.Cowboy{
			cbs[0], cbs[1],
		}, nil).Times(1)

	dbmocks.EXPECT().CreateGame(gomock.Any(), "basic").
		Return(uuid.New(), nil).Times(1)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("no rows in result set")).AnyTimes()

	b := &BasicGameEngine{
		Dao: dbmocks,
	}

	err := b.Run(ctx)
	time.Sleep(5 * time.Second)
	assert.NoError(t, err)

}

func Test_Run_Only_One_alive(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	cbs := []*db.Cowboy{
		{
			ID:        shooterId,
			CreatedAt: pgtype.Timestamptz{Time: time.Now().AddDate(0, 0, -1)},
			Name:      "John",
			Health:    1,
			Damage:    1,
		},
	}

	dbmocks.EXPECT().ListAliveCowboys(gomock.Any(), db.ListAliveCowboysParams{
		QuerySort:   "asc",
		QueryOffset: 0,
		QueryLimit:  5,
	}).
		Return([]*db.Cowboy{
			cbs[0],
		}, nil).Times(1)

	b := &BasicGameEngine{
		Dao: dbmocks,
	}

	err := b.Run(ctx)

	assert.Equal(t, errors.New("more than one cowboy is required to simulate this shooter, currently found : 1"), err)

}

func Test_StartGame_Shoot_And_Lose(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	// Set the http server mock
	server := prepareHttpServer(t)
	defer server.Close()

	endpointUlr := "http://" + server.Listener.Addr().String()

	// First iteration which allows the function to proceed
	cowboyData := getCowboysData(shooterId)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(&cowboyData[0], nil).Times(1)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(nil, errors.New("no rows in result set")).Times(1)

	b := &BasicGameEngine{
		Dao:     dbmocks,
		HttpUrl: endpointUlr,
	}
	err := b.StartGame(ctx, shooterId, gameId, "John")
	assert.NoError(t, err)
}

func Test_StartGame_Shoot_And_Win(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	// Set the http server mock
	server := prepareHttpServer(t)
	defer server.Close()

	endpointUlr := "http://" + server.Listener.Addr().String()

	cowboyData := getCowboysData(shooterId)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(&cowboyData[0], nil).Times(1)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(&cowboyData[1], nil).Times(1)

	// Setting the winner in DB
	dbmocks.EXPECT().UpdateGameWinner(gomock.Any(),
		db.UpdateGameWinnerParams{
			Winner:   cowboyData[1].Name,
			WinnerID: shooterId,
			GameID:   gameId,
		}).
		Return(nil).Times(1)

	b := &BasicGameEngine{
		Dao:     dbmocks,
		HttpUrl: endpointUlr,
	}

	err := b.StartGame(ctx, shooterId, gameId, "John")
	assert.NoError(t, err)
}

func Test_StartGame_Http_error(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	// Set the http server mock
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFailedDependency)
		_, err := w.Write([]byte(`{"error":"0"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	endpointUlr := "http://" + server.Listener.Addr().String()

	// First iteration which allows the function to proceed
	cowboyData := getCowboysData(shooterId)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(&cowboyData[0], nil).Times(1)

	//dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
	//	Return(nil, errors.New("no rows in result set")).Times(1)

	b := &BasicGameEngine{
		Dao:     dbmocks,
		HttpUrl: endpointUlr,
	}
	err := b.StartGame(ctx, shooterId, gameId, "John")
	assert.Error(t, err, "http request received statuscode = 424")
}

func Test_StartGame_Database_Disconnected(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(nil, errors.New("connect: connection refused")).Times(1)

	b := &BasicGameEngine{
		Dao: dbmocks,
	}
	err := b.StartGame(ctx, shooterId, gameId, "John")
	assert.Error(t, err, "connect: connection refused")
}

func setupMemDB() (*connection.TestConnection, error) {
	conn, err := connection.NewTestConnection()
	if err != nil {
		log.Fatalf("failed to create new cockroach connection: %v", err)
	}

	config, err := pgxpool.ParseConfig(conn.GetDsn())
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	config.MaxConns = 30 // I do this because im killing connections in my tests :P
	config.MinConns = 20

	dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("conencting to database with min %v, max %v pool", config.MinConns, config.MaxConns)
	}
	return conn, nil
}

func Test_ShootRandomCowboy(t *testing.T) {
	ctx := context.Background()

	conn, err := setupMemDB()

	// Migrate our db to mem
	migrationSource := &migrate.FileMigrationSource{Dir: "../../res/sql/migrations"}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("can't close connection for db migration")
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	n, err := migrate.Exec(conn.GetDb(), "postgres", migrationSource, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%v migrations run ", n)

	d := db.New(dbPool)

	//Clear default Db migrated cowboys
	_, err = dbPool.Exec(ctx, "Delete from cowboys;")
	if err != nil {
		t.Errorf("failed deleting existing migration data : %v", err)
	}

	_, err = dbPool.Exec(ctx, "INSERT INTO cowboys (id, name, health, damage) VALUES "+
		fmt.Sprintf("('%s', 'John2', 10, 1), ('%s','Bill2', 8, 2);", shooterId, "3bf49568-e955-4d2f-b78f-7df603571d31"))
	if err != nil {
		t.Errorf("failed while preparing example data : %v", err)
	}

	b := &BasicGameEngine{
		Dao: d,
		Db:  dbPool,
	}

	receiverHealth, err := b.ShootRandomCowboy(ctx, shooterId, gameId, "John2", 1)
	assert.NoError(t, err)
	assert.Equal(t, int32(7), receiverHealth)

}

func Test_ShootRandomCowboy_Last_Survivor(t *testing.T) {
	ctx := context.Background()

	conn, err := setupMemDB()

	// Migrate our db to mem
	migrationSource := &migrate.FileMigrationSource{Dir: "../../res/sql/migrations"}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("can't close connection for db migration")
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	n, err := migrate.Exec(conn.GetDb(), "postgres", migrationSource, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%v migrations run ", n)

	d := db.New(dbPool)

	//Clear default Db migrated cowboys
	_, err = dbPool.Exec(ctx, "Delete from cowboys;")
	if err != nil {
		t.Errorf("failed deleting existing migration data : %v", err)
	}

	_, err = dbPool.Exec(ctx, "INSERT INTO cowboys (id, name, health, damage) VALUES "+
		fmt.Sprintf("('%s', 'John2', 10, 1);", shooterId))
	if err != nil {
		t.Errorf("failed while preparing example data : %v", err)
	}

	b := &BasicGameEngine{
		Dao: d,
		Db:  dbPool,
	}

	receiverHealth, err := b.ShootRandomCowboy(ctx, shooterId, gameId, "John2", 1)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), receiverHealth)

}

func Test_NewBasicGameEngine(t *testing.T) {
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))
	urltest := "http://testing.com"
	os.Setenv("BASIC_GAME_URL", urltest)
	eng := NewBasicGameEngine(dbmocks, dbPool)

	assert.Equal(t, dbmocks, eng.Dao)
	assert.Equal(t, dbPool, eng.Db)
	assert.Equal(t, urltest, eng.HttpUrl)
	os.Clearenv()
}

func Test_NewBasicGameEngine_Url_Not_Set(t *testing.T) {

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))
	urlDefault := "http://localhost:9090"
	eng := NewBasicGameEngine(dbmocks, dbPool)

	assert.Equal(t, dbmocks, eng.Dao)
	assert.Equal(t, dbPool, eng.Db)
	assert.Equal(t, urlDefault, eng.HttpUrl)
}

func prepareHttpServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shoot" {
			t.Errorf("Expected to request '/shoot', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"receiver_health":"0"}`))
		assert.NoError(t, err)
	}))
}

func getCowboysData(shooterId uuid.UUID) []db.GetSingleAliveCowboyAndCountRow {
	// First iteration which allows the function to proceed
	return []db.GetSingleAliveCowboyAndCountRow{
		{
			ID:        shooterId,
			CreatedAt: pgtype.Timestamptz{Time: time.Now()},
			Name:      "John",
			Health:    1,
			Damage:    1,
			Available: 1,
		},
		{
			ID:        shooterId,
			CreatedAt: pgtype.Timestamptz{Time: time.Now()},
			Name:      "John",
			Health:    1,
			Damage:    1,
			Available: 0,
		},
	}
}

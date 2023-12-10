package intermediate_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/sumlookup/cowboys/dao/connection"
	"github.com/sumlookup/cowboys/dao/db"
	inter "github.com/sumlookup/cowboys/engine/intermediate"
	han "github.com/sumlookup/cowboys/handler"
	"github.com/sumlookup/cowboys/mocks"
	pb "github.com/sumlookup/cowboys/pb"
	"github.com/sumlookup/mini/service"
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
	testSvc   *service.Service
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
	testSvc = service.NewService("test", "grpc", "mdns")
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

	dbmocks.EXPECT().CreateGame(gomock.Any(), "intermediate").
		Return(uuid.New(), nil).Times(1)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("no rows in result set")).AnyTimes()
	b := &inter.IntermediateGameEngine{
		Dao:     dbmocks,
		Service: testSvc,
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

	b := &inter.IntermediateGameEngine{
		Dao: dbmocks,
	}

	err := b.Run(ctx)

	assert.Equal(t, errors.New("more than one cowboy is required to simulate this shooter, currently found : 1"), err)

}

func Test_StartGame_No_More_left(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(nil, errors.New("no rows in result set")).Times(1)
	os.Clearenv()
	os.Setenv("SELECTOR", "registry")
	b := &inter.IntermediateGameEngine{
		Dao:     dbmocks,
		Service: testSvc,
	}
	mockEngine := mocks.NewMockCowboysEngine(gomock.NewController(t))

	mockEngine.EXPECT().ShootRandomCowboy(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	testService := service.NewService("cowboys", "grpc", "mdns")
	pb.RegisterCowboysServiceServer(testService.Server(), &han.CowboysService{
		Dao:    dbmocks,
		Engine: mockEngine,
	})

	go func() {
		err := testService.Run()
		if err != nil {
			t.Errorf("failed running test resvice")
		}
	}()
	defer testService.Server().Stop()
	defer testService.Close()

	err := b.StartGame(ctx, shooterId, gameId, "John")
	assert.NoError(t, err)
}

func Test_StartGame_Winner(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	// First iteration which allows the function to proceed
	cowboyData := getCowboysData(shooterId)

	dbmocks.EXPECT().GetSingleAliveCowboyAndCount(gomock.Any(), shooterId).
		Return(&cowboyData[1], nil).Times(1)

	dbmocks.EXPECT().UpdateGameWinner(gomock.Any(),
		db.UpdateGameWinnerParams{
			Winner:   cowboyData[1].Name,
			WinnerID: shooterId,
			GameID:   gameId,
		}).
		Return(nil).Times(1)

	os.Clearenv()
	os.Setenv("SELECTOR", "registry")
	b := &inter.IntermediateGameEngine{
		Dao:     dbmocks,
		Service: testSvc,
	}
	mockEngine := mocks.NewMockCowboysEngine(gomock.NewController(t))

	mockEngine.EXPECT().ShootRandomCowboy(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	testService := service.NewService("cowboys", "grpc", "mdns")
	pb.RegisterCowboysServiceServer(testService.Server(), &han.CowboysService{
		Dao:    dbmocks,
		Engine: mockEngine,
	})

	go func() {
		err := testService.Run()
		if err != nil {
			t.Errorf("failed running test resvice")
		}
	}()
	defer testService.Server().Stop()
	defer testService.Close()

	err := b.StartGame(ctx, shooterId, gameId, "John")
	assert.NoError(t, err)
}

func Test_NewIntermediateGameEngine(t *testing.T) {
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))
	eng := inter.NewIntermediateGameEngine(dbmocks, dbPool, testSvc)

	assert.Equal(t, dbmocks, eng.Dao)
	assert.Equal(t, dbPool, eng.Db)
	assert.Equal(t, testSvc, eng.Service)
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
		t.Fatal(err)
	}

	n, err := migrate.Exec(conn.GetDb(), "postgres", migrationSource, migrate.Up)
	if err != nil {
		t.Fatal(err)
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

	b := &inter.IntermediateGameEngine{
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

	b := &inter.IntermediateGameEngine{
		Dao: d,
		Db:  dbPool,
	}

	receiverHealth, err := b.ShootRandomCowboy(ctx, shooterId, gameId, "John2", 1)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), receiverHealth)

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

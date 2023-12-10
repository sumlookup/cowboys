package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/sumlookup/cowboys/dao/connection"
	daodb "github.com/sumlookup/cowboys/dao/db"
	"github.com/sumlookup/cowboys/engine"
	"github.com/sumlookup/cowboys/gateway"
	han "github.com/sumlookup/cowboys/handler"
	pb "github.com/sumlookup/cowboys/pb"
	"github.com/sumlookup/mini/service"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
	"os"
)

func main() {

	svc := GetServiceName()
	app := service.NewService(svc, GetTransport(), GetRegistry())

	// make this available globally
	var s *http.Server

	// end of function, shutdown the server
	defer func() {
		if s != nil {
			if err := s.Shutdown(context.Background()); err != nil {
				log.Errorf("Failed to shutdown http gateway server: %v", err)
				return
			}
		}
	}()

	// start the http server
	go func() {

		nc := app.Client(GetSelector()).Connect(svc)
		s = gateway.NewGwServer(context.Background(), nc, GetPort())

		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.Infof("Serving http server on http://0.0.0.0%v", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Error(err)
			// any problems, shut down the app
			app.Close()
			return
		}
	}()
	url := os.Getenv("DB_URL")
	log.Infof("DB_URL - %v", url)
	dbg, err := connection.NewConnection(url)
	if err != nil {
		log.Fatalf("failed to establish db connection please set DB_URL env - %v", err)
	}

	pool := dbg.GetDbPool()
	db := daodb.New(pool)
	e := engine.NewEngine(GetGameMode(), db, pool, app)

	h := &han.CowboysService{
		Engine: e,
		Dao:    db,
		Db:     pool,
	}

	//RUN DB migrations
	go func() {
		var err error
		dbm, err := sql.Open("postgres", url)
		defer func() {
			err = dbm.Close()
			if err != nil {
				log.Errorf("can't close connection for db migration")
			}
		}()
		if err != nil {
			log.Fatalf("sql failed to establish connection with DB, is DB running?,  err : %v", err)
		}

		migrationSource := &migrate.FileMigrationSource{Dir: "./res/sql/migrations"}

		n, err := migrate.Exec(dbm, "postgres", migrationSource, migrate.Up)
		if err != nil {
			log.Fatalf("sql migration failed,  err : %v", err)
		}
		log.Infof("%v migrations run ", n)
	}()

	app.AddHandler(h)

	// Register the validator server
	pb.RegisterCowboysServiceServer(app.Server(), h)

	// Register the healthcheck server
	healthpb.RegisterHealthServer(app.Server(), h)

	err = app.Run()
	if err != nil {
		// any problems, log the error
		log.Error(err)
	}

	app.Close()
}

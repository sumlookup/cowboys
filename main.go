package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	con "github.com/sumlookup/cowboys/dao/connection"
	daodb "github.com/sumlookup/cowboys/dao/db"
	"github.com/sumlookup/cowboys/engine"
	gateway "github.com/sumlookup/cowboys/gateway"
	han "github.com/sumlookup/cowboys/handler"
	pb "github.com/sumlookup/cowboys/pb"
	"github.com/sumlookup/mini/service"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
	"os"
)

//
//Requirements
//Input
//We have a set of cowboys.
//Each cowboy has a unique name, health points and damage points.
//Cowboys list must be stored in persistent storage (File, Database etc).
//Each cowboy should run in it’s own isolated process, workload or replica.
//All communication between cowboys should happen via your preferred networking solution (TCP, gRPC, HTTP, MQ etc). Cowboys encounter starts at the same time in parallel. Each cowboys selects random target and shoots.
//Subtract shooter damage points from target health points.
//If target cowboy health points are 0 or lower, then target is dead.
//Cowboys don’t shoot themselves and don’t shoot dead cowboys.
//After the shot shooter sleeps for 1 second.
//Last standing cowboy is the winner.
//Outcome of the task is to print log of every action and winner should log that he won.
//Kubernetes, Docker-compose or any other container orchestration solution is preferred, but optional for final deployment manifests. Provide startup and usage instructions in Readme.MD

const (
	DEFAULT_PORT = "9091"
	SERVICE_NAME = "cowboys-svc"
	TRANSPORT    = "grpc"
	REGISTRY     = "mdns"
	SELECTOR     = "registry"

	MODE = "basic"
)

func main() {

	app := service.NewService(SERVICE_NAME, GetTransport(), GetRegistry())

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

		nc := app.Client(GetSelector()).Connect(SERVICE_NAME)
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
	log.Infof("db url - %v", url)
	dbg, err := con.NewConnection(url)
	if err != nil {
		log.Fatal("failed to establish db connection please set DB_URL env - %v", err)
	}

	pool := dbg.GetDbPool()
	db := daodb.New(pool)
	e := engine.NewEngine(GetGameMode(), db, pool, app)

	h := &han.CowboysService{
		Engine: e,
		Dao:    db,
		Db:     pool,
	}

	// RUN DB migrations
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
			log.Fatalf("sql fialed to establish connection with DB, is DB running?,  err : %v", err)
		}

		migrationSource := &migrate.FileMigrationSource{Dir: "res/sql/migrations"}

		n, err := migrate.Exec(dbm, "postgres", migrationSource, migrate.Up)
		if err != nil {
			log.Fatalf("sql migration failed,  err : %v", err)
		}
		log.Infof("%v migrations run ", n)
	}()

	app.AddHandler(h)

	// Register the validator server
	pb.RegisterCowboysServiceServer(app.Server(), h)

	// Register the healtcheck server
	healthpb.RegisterHealthServer(app.Server(), h)

	err = app.Run()
	if err != nil {
		// any problems, log the error
		log.Error(err)
	}

	app.Close()
}

// todo: move bellow funcs to other place

func GetPort() string {
	port := DEFAULT_PORT
	value := os.Getenv("HTTP_PORT")
	if len(value) == 0 {
		return port
	}
	return value
}

func GetSelector() string {
	selector := SELECTOR
	value := os.Getenv("SELECTOR")
	if len(value) == 0 {
		return selector
	}
	return value
}

func GetTransport() string {
	transport := TRANSPORT
	value := os.Getenv("TRANSPORT")
	if len(value) == 0 {
		return transport
	}
	return value
}

func GetRegistry() string {
	registry := REGISTRY
	value := os.Getenv("REGISTRY")
	if len(value) == 0 {
		return registry
	}
	return value
}

func GetGameMode() string {
	registry := MODE
	value := os.Getenv("MODE")
	if len(value) == 0 {
		return registry
	}
	return value
}

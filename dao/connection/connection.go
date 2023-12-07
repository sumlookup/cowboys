package connection

import (
	"context"
	pgxv5 "github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

// https://www.cockroachlabs.com/docs/v22.1/connection-pooling?filters=go#sizing-connection-pools
var poolMaxConn = int32(runtime.NumCPU() * 4)

// Cockroach Labs recommends setting the maximum number of idle connections to the maximum pool size.
// While this uses more memory, it allows many connections when concurrency is high without having to create a new connection for every new operation.
var poolMinConn = int32(runtime.NumCPU() * 4)

type Connection struct {
	dbPool *pgxv5.Pool
}

var MaxAttempts int64 = 5

func NewConnection(connString string) (*Connection, error) {
	u, err := url.Parse(connString)
	if err != nil {
		return nil, err
	}
	config, err := pgxv5.ParseConfig(connString)
	if err != nil {
		log.Errorf("failed while attempting to connect with pgx - %v", err)
		return nil, err
	}
	param, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	if pmc, ok := param["pool_max_conns"]; ok {
		maxConns, err := strconv.ParseInt(pmc[0], 10, 32)
		if err != nil {
			return nil, err
		}
		config.MaxConns = int32(maxConns)
	} else {
		config.MaxConns = poolMaxConn
	}

	if pmc, ok := param["pool_min_conns"]; ok {
		minConns, err := strconv.ParseInt(pmc[0], 10, 32)
		if err != nil {
			return nil, err
		}
		config.MinConns = int32(minConns)
	} else {
		config.MinConns = poolMinConn
	}

	ticker := time.NewTicker(1 * time.Second)
	var i int64 = 0

	maxAttempts, err := strconv.ParseInt(os.Getenv("DB_MAX_ATTEMPTS"), 10, 32)
	if err != nil {
		maxAttempts = MaxAttempts
	}

	for range ticker.C {
		if i >= maxAttempts {
			break
		}

		dbPool, pgxErr := pgxv5.NewWithConfig(context.Background(), config)
		if pgxErr == nil {
			log.Infof("conencting to database with min %v, max %v pool", config.MinConns, config.MaxConns)
			return &Connection{dbPool: dbPool}, nil
		} else {
			i++
			log.Error("error connecting to the database, trying again: ", pgxErr)
		}
	}
	log.Fatal("can't connect to the database: ", err)
	return nil, nil
}

func (c *Connection) Close() {
	c.dbPool.Close()
}

func (c *Connection) GetDbPool() *pgxv5.Pool {
	return c.dbPool
}

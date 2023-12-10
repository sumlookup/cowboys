package connection

import (
	//"context"
	"database/sql"
	//"fmt"

	testserverCR "github.com/cockroachdb/cockroach-go/v2/testserver"
	log "github.com/sirupsen/logrus"
	//"gorm.io/driver/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestConnection struct {
	Db      *sql.DB
	PgxPool *pgxpool.Pool
	ts      testserverCR.TestServer
	Dsn     string
}

// NewConnection creates new DB connection
func NewTestConnection() (*TestConnection, error) {
	ts, err := testserverCR.NewTestServer(testserverCR.CustomVersionOpt("22.1.0"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", ts.PGURL().String())

	if err != nil {
		log.Fatalf("failed to create new DB connection: %v", err)
	}

	return &TestConnection{ts: ts, Db: db, Dsn: ts.PGURL().String()}, nil
}

func (c *TestConnection) GetDb() *sql.DB {
	return c.Db
}

func (c *TestConnection) GetDsn() string {
	return c.Dsn
}

func (c *TestConnection) Close() error {
	log.Info("Stopping DB Testserver Postgres")
	c.ts.Stop()

	return nil
}

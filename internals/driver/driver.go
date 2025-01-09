package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

// Maximum number of db connection at a given instance
const maxOpenConn = 10
const maxIdleConn = 5
const maxDBLifetime = 5 * time.Minute

// ConnectionSQL Create database pool for Postgres
func ConnectionSQL(dsn string) (*DB, error) {
	dbPool, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	dbPool.SetMaxOpenConns(maxOpenConn)
	dbPool.SetMaxIdleConns(maxIdleConn)
	dbPool.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = dbPool

	err = testDB(dbPool)
	if err != nil {
		panic(err)
	}
	return dbConn, err
}

// testDB tries to ping the database
func testDB(dbPool *sql.DB) error {
	err := dbPool.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase create new database for the connection
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err

	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}

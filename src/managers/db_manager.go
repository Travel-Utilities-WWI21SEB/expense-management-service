package managers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
)

type DatabaseMgr interface {
	ExecuteStatement(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecuteQuery(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	ExecuteQueryRow(ctx context.Context, query string, args ...any) pgx.Row
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CheckIfExists(ctx context.Context, query string, args ...any) (bool, error)
}

type DatabaseManager struct {
	Connection *pgxpool.Pool
}

// BeginTx starts a transaction and returns a transaction object, this is helpful because we can pass the transaction object to the repositories
// and they can use it to execute multiple statements in a single transaction
// TODO: Use this in the repositories
func (dm *DatabaseManager) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := dm.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (dm *DatabaseManager) ExecuteStatement(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	result, err := dm.Connection.Exec(ctx, query, args...)
	return result, err
}

func (dm *DatabaseManager) ExecuteQuery(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	rows, err := dm.Connection.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (dm *DatabaseManager) ExecuteQueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	row := dm.Connection.QueryRow(ctx, query, args...)
	return row
}

func (dm *DatabaseManager) CheckIfExists(ctx context.Context, query string, args ...any) (bool, error) {
	var count int
	err := dm.Connection.QueryRow(ctx, query, args...).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InitializeDatabaseConnection(ctx context.Context) (*pgxpool.Pool, error) {
	// Check if environment variables are set
	var environment = os.Getenv("ENVIRONMENT")

	var (
		dbHost     = os.Getenv(fmt.Sprintf("%s_DB_HOST", environment))
		dbPort     = os.Getenv(fmt.Sprintf("%s_DB_PORT", environment))
		dbUser     = os.Getenv(fmt.Sprintf("%s_DB_USER", environment))
		dbPassword = os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", environment))
		dbName     = os.Getenv(fmt.Sprintf("%s_DB_NAME", environment))
	)

	if utils.ContainsEmptyString(dbHost, dbPort, dbUser, dbPassword, dbName) {
		return nil, fmt.Errorf("error initializing database connection: environment variables not set")
	}

	// CONNECT TO DATABASE
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	config, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	config.MinConns = 5
	config.MaxConns = 30
	config.MaxConnLifetime = time.Minute * 30
	config.MaxConnIdleTime = time.Minute * 5
	config.HealthCheckPeriod = time.Minute * 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	// pool, err := pgxpool.New(ctx, psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error creating new pool: %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return pool, nil
}

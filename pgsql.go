package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type handler struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) handler {
	return handler{db}
}

func Connect() *pgxpool.Pool {
	connInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	pgxConfig, err := pgxpool.ParseConfig(connInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pgxConfig.MaxConns = 5
	pgxConfig.MaxConnLifetime = 5 * time.Minute
	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.HealthCheckPeriod = time.Minute
	pgxConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to db!")
	return db
}

func CloseConnection(db *pgxpool.Pool) {
	defer db.Close()
}

package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() *DB {
	url := os.Getenv("DATABASE_URL")
	
	if url == "" {
		log.Fatal("DATABASE_URL env variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	//Test the conection with a ping
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database %v", err)
	}

	log.Println("Successfully connected to database")

	return &DB{Pool: pool}
}

func (db *DB) Close() {
	db.Pool.Close()
}

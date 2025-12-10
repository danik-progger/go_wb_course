package db

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"salesTracker/db/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

type DB struct {
	ctx context.Context
	q   *repo.Queries
}

func InitDB() (*DB, error) {
	ctx := context.Background()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// If .env file is not found, continue with environment variables
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	// Get database connection parameters from environment
	user := getEnv("DB_USER", "pqgotest")
	dbname := getEnv("DB_NAME", "pqgotest")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	sslmode := getEnv("DB_SSLMODE", "verify-full")

	// Build connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user, password, dbname, host, port, sslmode)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	queries := repo.New(conn)

	return &DB{
		ctx: ctx,
		q:   queries,
	}, nil
}

// getEnv retrieves environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (d *DB) GetAll() ([]repo.Sale, error) {
	return d.q.GetAll(d.ctx)
}

func (d *DB) GetById(id int64) (repo.Sale, error) {
	return d.q.GetById(d.ctx, id)
}

func (d *DB) GetByCategory(category string) ([]repo.Sale, error) {
	categoryText := pgtype.Text{String: category, Valid: true}
	return d.q.GetByCategory(d.ctx, categoryText)
}

func (d *DB) CreateSale(category string, date pgtype.Date, amount pgtype.Numeric) (repo.Sale, error) {
	params := repo.CreateSaleParams{
		Category: pgtype.Text{String: category, Valid: true},
		Date:     date,
		Amount:   amount,
	}
	return d.q.CreateSale(d.ctx, params)
}

func (d *DB) DeleteItem(id int64) error {
	return d.q.DeleteItem(d.ctx, id)
}

func (d *DB) GetAverage(date1, date2 pgtype.Date) (float64, error) {
	params := repo.GetAverageParams{
		Date:   date1,
		Date_2: date2,
	}
	return d.q.GetAverage(d.ctx, params)
}

func (d *DB) GetDateInterval(dateStart, dateEnd pgtype.Date) ([]repo.Sale, error) {
	params := repo.GetDateIntervalParams{
		Date:   dateStart,
		Date_2: dateEnd,
	}
	return d.q.GetDateInterval(d.ctx, params)
}

func (d *DB) GetMedian(date1, date2 pgtype.Date) (float64, error) {
	params := repo.GetMedianParams{
		Date:   date1,
		Date_2: date2,
	}
	return d.q.GetMedian(d.ctx, params)
}

func (d *DB) GetPercentile(percentileCont float64, date1, date2 pgtype.Date) (float64, error) {
	params := repo.GetPercentileParams{
		PercentileCont: percentileCont,
		Date:           date1,
		Date_2:         date2,
	}
	return d.q.GetPercentile(d.ctx, params)
}

func (d *DB) GetSum(date1, date2 pgtype.Date) (int64, error) {
	params := repo.GetSumParams{
		Date:   date1,
		Date_2: date2,
	}
	return d.q.GetSum(d.ctx, params)
}

func (d *DB) UpdateItem(id int64, category string, date pgtype.Date, amount pgtype.Numeric) (repo.Sale, error) {
	params := repo.UpdateItemParams{
		ID:       id,
		Category: pgtype.Text{String: category, Valid: true},
		Date:     date,
		Amount:   amount,
	}
	return d.q.UpdateItem(d.ctx, params)
}

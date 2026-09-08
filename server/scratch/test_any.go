package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgres://postgres:password123@localhost:5432/tnnow_dev?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("DB error: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	ids := []string{"3b79f87c-39b2-487d-a778-2c7fa488ffcd"}

	// Test WHERE id = ANY($1) with []string
	tag, err := pool.Exec(ctx, "UPDATE content SET updated_at = NOW() WHERE id = ANY($1)", ids)
	fmt.Printf("Test 1 [WHERE id = ANY($1)]: tag=%v, err=%v\n", tag, err)

	// Test WHERE id::text = ANY($1)
	tag2, err2 := pool.Exec(ctx, "UPDATE content SET updated_at = NOW() WHERE id::text = ANY($1)", ids)
	fmt.Printf("Test 2 [WHERE id::text = ANY($1)]: tag=%v, err=%v\n", tag2, err2)
}

// Command migrate applies or rolls back database migrations by hand.
//
//	go run ./cmd/migrate up        apply all pending migrations
//	go run ./cmd/migrate down      roll back the last migration
//	go run ./cmd/migrate down 4    roll back the last 4 migrations
//	go run ./cmd/migrate version   show the current version
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"meditrack/internal/db"
)

func main() {
	_ = godotenv.Load()
	url := os.Getenv("DATABASE_URL")
	if url == "" || len(os.Args) < 2 {
		fmt.Println("usage: DATABASE_URL=... go run ./cmd/migrate up|down [n]|version")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "up":
		err = db.MigrateUp(url)
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			if steps, err = strconv.Atoi(os.Args[2]); err != nil || steps < 1 {
				fmt.Println("down needs a positive number of steps")
				os.Exit(1)
			}
		}
		err = db.MigrateDown(url, steps)
	case "version":
		v, dirty, verr := db.MigrateVersion(url)
		err = verr
		if err == nil {
			fmt.Printf("version %d (dirty: %v)\n", v, dirty)
		}
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Println("done")
}

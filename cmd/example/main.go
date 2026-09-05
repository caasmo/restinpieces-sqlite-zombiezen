package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caasmo/restinpieces"
	"github.com/caasmo/restinpieces-sqlite-zombiezen"
	"github.com/caasmo/restinpieces-sqlite-zombiezen/zombiezen"
)

func main() {
	dbPath := flag.String("db", "", "Path to the SQLite database file (required)")
	ageKeyPath := flag.String("age-key", "", "Path to the age identity (private key) file (required)")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s -db <database-path> -age-key <identity-file-path>\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Start the restinpieces application server.\n\n")
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *dbPath == "" || *ageKeyPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// --- Create the Database Pool ---
	// Use the driver's pool constructor with restinpieces-suitable defaults
	// (WAL mode, busy_timeout). Main owns the pool lifecycle.
	pool, err := zombiezen.NewPool(*dbPath)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		slog.Info("Closing database pool...")
		closeErr := pool.Close()
		if closeErr != nil {
			slog.Error("Error closing database pool", "error", closeErr)
		}
	}()

	// --- Initialize the Application ---
	_, srv, err := restinpieces.New(
		restinpieces.WithAgeKeyPath(*ageKeyPath),
		sqlitezombiezen.WithDbZombiezen(pool),
	)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	srv.Run()

	slog.Info("Server shut down gracefully.")
}

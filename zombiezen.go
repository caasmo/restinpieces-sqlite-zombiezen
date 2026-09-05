package sqlitezombiezen

import (
	"fmt"

	"github.com/caasmo/restinpieces"
	"github.com/caasmo/restinpieces-sqlite-zombiezen/zombiezen"
	"zombiezen.com/go/sqlite/sqlitex"
)

// WithDbZombiezen configures the App to use the Zombiezen SQLite implementation with an existing pool.
func WithDbZombiezen(pool *sqlitex.Pool) restinpieces.Option {
	dbInstance, err := zombiezen.New(pool)
	if err != nil {
		// Panic is reasonable here as it indicates a fundamental setup error.
		panic(fmt.Sprintf("failed to initialize zombiezen DB with existing pool: %v", err))
	}

	return restinpieces.WithDbApp(dbInstance)
}

// If your application interacts directly with the database alongside restinpieces,
// it's crucial to use a *single shared pool* to prevent database locking issues (SQLITE_BUSY errors).
// Create the pool with zombiezen.NewPool (WAL mode and busy_timeout defaults
// suitable for restinpieces) and pass it to both restinpieces (via WithDbZombiezen)
// and your own application's database access layer.

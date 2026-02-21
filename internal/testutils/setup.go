package testutils

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// Initializes an in-memory SQLite database, runs migrations,
// and returns a repository.Database instance.
func SetupTestDB(t *testing.T) repository.Database {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(t, err, "Failed to open in-memory database")

	t.Cleanup(func() {
		db.Close()
	})

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	require.NoError(t, err, "Failed to create migration driver")

	// Locate migrations directory
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "../../..")
	migrationsPath := "file://" + filepath.Join(projectRoot, "server/migrations")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3",
		driver,
	)
	require.NoError(t, err, "Failed to initialize migrations")

	err = m.Up()
	if err != migrate.ErrNoChange {
		require.NoError(t, err, "Failed to run migrations")
	}

	return repository.NewDatabase(db)
}

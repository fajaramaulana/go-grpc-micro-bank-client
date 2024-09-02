package dbmigration

import (
	"database/sql"

	"github.com/fajaramaulana/go-grpc-micro-bank-server/util"
	migrate "github.com/golang-migrate/migrate/v4"

	"github.com/chilts/sid"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

// Migrate migrates the database using the provided SQL connection.
// It performs the following steps:
// 1. Generates a unique session ID.
// 2. Logs the start of the migration process.
// 3. Sets up the database driver using the provided SQL connection.
// 4. Creates a new migration instance with the specified migration files directory and database driver.
// 5. Executes the "up" migration.
// 6. Logs the result of the migration process.
//
// If an error occurs during any of the steps, the function logs the error and terminates the migration process.
// The function returns no values.
func Migrate(conn *sql.DB) {
	sidString := sid.Id()
	logStart := util.LogRequest("", "Migrate-"+sidString, "dbmigration - Migrate")
	log.Info().Msg(logStart)

	dbDriver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		logErr := util.LogError(err.Error(), "Migrate-"+sidString, "dbmigration - Migrate - postgres.WithInstance")
		log.Fatal().Msg(logErr)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../db/migrations", "postgres", dbDriver)

	if err != nil {
		logErr := util.LogError(err.Error(), "Migrate-"+sidString, "dbmigration - Migrate - migrate.NewWithDatabaseInstance")
		log.Fatal().Msg(logErr)
	}

	// if err := m.Down(); err != nil {
	// 	log.Println("Database migration (down) failed :", err)
	// }

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Info().Msg("Database migration (up) no change")
		} else {
			logErr := util.LogError(err.Error(), "Migrate-"+sidString, "dbmigration - Migrate - m.Up()")
			log.Fatal().Msg(logErr)
		}
	} else {
		logStop := util.LogResponse("Database Migration Done", "Migrate-"+sidString, "dbmigration - Migrate")
		log.Info().Msg(logStop)
	}
}

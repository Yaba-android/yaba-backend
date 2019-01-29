package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
	"github.com/nasrat_v/maktaba-android-mvp/src/controllers/ebook"
	"github.com/nasrat_v/maktaba-android-mvp/src/controllers/library"
	"github.com/nasrat_v/maktaba-android-mvp/src/controllers/user"
)

var dbHandler *pgx.Conn

// initDatabaseConnection initialize connection with the database
func initDatabaseConnection() {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		fmt.Println("Unable to parse Postgre environment: ", err)
		os.Exit(1)
	}
	dbHandler, err = pgx.Connect(config)
	if err != nil {
		fmt.Println("Connection to database failed: ", err)
		os.Exit(1)
	}
}

// InitDbHandlerForControllers set the database handler to all controllers
func InitDbHandlerForControllers() {
	initDatabaseConnection()
	ebook.DbHandler = dbHandler
	user.DbHandler = dbHandler
	library.DbHandler = dbHandler
}

package commons

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func GetConnection() (db *sql.DB) {

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}

	// Get a database handle.

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Print(err)
		panic(err.Error())
	}
	return db
}

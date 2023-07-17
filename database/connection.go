package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// GetDB returns a connection to the MySQL database
func GetDB() *sql.DB {
	db, err := sql.Open("mysql", "root:2309@/e_voting")
	if err != nil {
		panic(err)
	}
	return db
}

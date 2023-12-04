package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// GetDB returns a connection to the MySQL database
//replace MYSQL_PASS and MYSQL_USER with your database username and password
func GetDB() *sql.DB {
	db, err := sql.Open("mysql", "MYSQL_USER:MYSQL_PASS@/e_voting")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	} else{
		log.Println("Connected to the database successfully.")

	}
	return db
}

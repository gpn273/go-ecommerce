package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"os"
)

func DatabaseContext() *sql.DB {
	dataSourceName := DB_CONNECTION_STRING
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	connectionCheck := db.Ping()
	if connectionCheck != nil {
		fmt.Println("Application is unable to start as it cannot establish a connection to the MySQL database/server")
		os.Exit(1)
	}

	return db
}

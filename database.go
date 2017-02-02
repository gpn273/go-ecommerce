package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

func DatabaseContext() *sql.DB {
	dataSourceName := "root:password@/baltic"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	return db
}

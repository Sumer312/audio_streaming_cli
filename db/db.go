package db

import (
	"database/sql"
	"log"
)

func InitializeDB() {
	db, err := sql.Open("sqlite3", "./songsdb.sqlite")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
}

func AddToDB(song string) bool{

	return false
}

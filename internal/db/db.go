package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return DB.Exec(query, args...)
}

func init() {
	var err error
	if err != nil {
		panic(err)
	}
	DB, err = sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY, title TEXT, content TEXT)")
	if err != nil {
		panic(err)
	}
}

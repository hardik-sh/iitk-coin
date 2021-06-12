package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func search(rollno int, password string) bool {
	db, _ := sql.Open("sqlite3", "mydb.db")
	query := "SELECT EXISTS(SELECT 1 FROM Users WHERE rollno=? AND password=? LIMIT 1)"
	statement, _ := db.Prepare(query)
	bval, _ := statement.Exec(rollno, password)
	if bval != nil {
		return true
	} else {
		return false
	}
}
func addStudent(u User) {
	db, _ := sql.Open("sqlite3", "mydb.db")
	query := "CREATE TABLE IF NOT EXISTS Users (rollno INTEGER, password TEXT,username TEXT,UNIQUE(rollno))"
	statement, _ := db.Prepare(query)
	statement.Exec()
	query = "INSERT OR IGNORE INTO Users(rollno, password,username) VALUES (?, ?,?)"
	tmp, _ := db.Prepare(query)
	tmp.Exec(u.Rollno, u.Password, u.Username)
}

package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//to be run only when initialising the coin
func initialRoles() error {
	var fileName = "sqlite.db"
	var db, err = sql.Open("sqlite3", fileName)
	addQuery := "INSERT OR IGNORE INTO Users(rollno, password,username,coins,role) VALUES (?, ?,?,?,?)"
	pass, _ := bcrypt.GenerateFromPassword([]byte("iitkcoin"), 16)
	password := string(pass)
	db.Exec(addQuery, 1, password, "gensec", 0, 3)
	db.Exec(addQuery, 2, password, "AH1", 0, 3)
	db.Exec(addQuery, 3, password, "AH2", 0, 3)
	db.Exec(addQuery, 4, password, "AH3", 0, 3)
	db.Exec(addQuery, 5, password, "sysadmin", 0, 3)
	return err
}

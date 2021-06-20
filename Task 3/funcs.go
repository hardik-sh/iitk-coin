package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func transferhandler(w http.ResponseWriter, r *http.Request) {
	fileName := "sqlite.db"
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var user transfer
	err = decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("BEGIN TRANSACTION; UPDATE Users SET coins=coins-? WHERE rollno=?; UPDATE Users SET coins=coins+? WHERE rollno=?;COMMIT;", user.Coins, user.Sender, user.Coins, user.Receiver)
	if err != nil {
		fmt.Println(err)
	}
}
func addhandler(w http.ResponseWriter, r *http.Request) {
	fileName := "sqlite.db"
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var user add
	err = decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO Users(rollno, coins) VALUES(?,?);", user.User, user.Coins)
	if err.Error() == "SQLITE_CONSTRAINT_PRIMARYKEY" {
		_, err = db.Exec("UPDATE Users SET coins=coins+? WHERE rollno=?;", user.Coins, user.User)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
func getcoinshandler(w http.ResponseWriter, r *http.Request) {
	fileName := "sqlite.db"
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var user get
	err = decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec("SEARCH coins FROM Users WHERE rollno=? ;", user.User)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

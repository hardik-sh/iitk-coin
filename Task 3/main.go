package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

type transfer struct {
	Coins    int `json:"Coins"`
	Receiver int `json:"Reciever"`
	Sender   int `json:"Sender"`
}
type add struct {
	Coins int `json:"Coins"`
	User  int `json:"User"`
}
type get struct {
	User int `json:"User"`
}

func main() {
	database, _ := sql.Open("sqlite3", "mydb.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Users (rollno INTEGER , coins INTEGER DEFAULT 0,CHECK(coins>=0),UNIQUE(rollno))")
	statement.Exec()
	fmt.Printf("Starting server at port 8080\n")
	http.HandleFunc("/", reqh)
	http.ListenAndServe(":8080", nil)
}
func reqh(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.URL.Path == "/transfer" {
			transferhandler(w, r)
		} else if r.URL.Path == "/add" {
			addhandler(w, r)
		} else {
			http.Error(w, "Request not supported", http.StatusNotFound)
			return
		}
	} else if r.Method == "GET" {
		if r.URL.Path == "/getcoins" {
			getcoinshandler(w, r)
		} else {
			http.Error(w, "Request not supported", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
}

package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
	Username string `json:"username"`
}
type LoginField struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}
type Claims struct {
	Rollno int `json:"rollno"`
	jwt.StandardClaims
}
type Token struct {
	Jwt string `json:"jwt"`
}

var jwtKey = []byte("random_generator")

func main() {
	database, _ := sql.Open("sqlite3", "mydb.db")
	statement, error := database.Prepare("CREATE TABLE IF NOT EXISTS Users(rollno INTEGER,name TEXT)")
	if error != nil {
		fmt.Println(error)
	}
	statement.Exec()
	fmt.Printf("Starting server at port 8080\n")
	http.HandleFunc("/", reqh)
	http.ListenAndServe(":8080", nil)
}
func reqh(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}
	if r.URL.Path == "/login" {
		loginHandler(w, r)
	} else if r.URL.Path == "/signup" {
		signupHandler(w, r)
	} else if r.URL.Path == "/secretpage" {
		secrethandler(w, r)
	} else {
		http.Error(w, "Request not supported", http.StatusNotFound)
		return
	}
}

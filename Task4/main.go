package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/add", add)
	http.HandleFunc("/send", send)
	http.HandleFunc("/balance", balance)
	http.HandleFunc("/changePassword", changePass)
	http.HandleFunc("/changeRole", changeRole)
	startDB()
	initialRoles()
	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)
}

//Maximum cap on coins is set to be 1000000
//Roles
//0=>Normal
//1=>core team excet AH and gensec
//2=>AH and gensec and sysadmin
//3=>admin accounts

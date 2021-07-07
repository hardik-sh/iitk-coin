package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var fileName = "sqlite.db"
var db, err = sql.Open("sqlite3", fileName)

func balance(res http.ResponseWriter, req *http.Request) {
	var token Token
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&token)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		claims := &jwtClaims{}
		tkn, err := jwt.ParseWithClaims(token.Jwt, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		balance := GetBalance(claims.Rollno)
		var response Coins
		response.Balance = balance
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(response)
		addLog(fmt.Sprintf("%s%d%s%d", "Balance checked by user: ", claims.Rollno, "\t Balance: ", balance))

	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(fmt.Sprintln("Only Post Method allowed")))
	}
}
func add(res http.ResponseWriter, req *http.Request) {
	var user addRoute
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		claims := &jwtClaims{}
		tkn, err := jwt.ParseWithClaims(user.Jwt, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims.Role != 3 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		role := GetRole(user.Rollno)
		if role == 2 || role == 3 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = AddCoin(user)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Some error occured please try again later"))

			return
		}
		res.WriteHeader(http.StatusAccepted)
		addLog(fmt.Sprintf("%s %d %s %d", "Money added to the account of user : ", user.Rollno, "by core team member: ", claims.Rollno))
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(fmt.Sprintln("Only Post Method allowed")))
	}
}
func send(res http.ResponseWriter, req *http.Request) {
	var user sendRoute
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		claims := &jwtClaims{}
		tkn, err := jwt.ParseWithClaims(user.SenderJWT, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			fmt.Println(err)
			if err == jwt.ErrSignatureInvalid {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		recieverRole := GetRole(user.Reciever)
		fmt.Println(recieverRole)
		if recieverRole == -1 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if recieverRole != 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		SendCoin(user, claims.Rollno, claims.Rollno/10000 == user.Reciever/10000)
		addLog(fmt.Sprintf("%d %s %d %s %d", claims.Rollno, " Sent ", user.Coins, " coins to ", user.Reciever))
		res.WriteHeader(http.StatusAccepted)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(("Only POST Method allowed")))
	}
}
func login(res http.ResponseWriter, req *http.Request) {
	var user loginRoute
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		if user.Password == "" || user.Rollno == 0 {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}

		role := SearchStudent(user)
		if role == -1 {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("User doesn't exist"))
			return
		}

		expiry := time.Now().Add(10 * time.Minute)
		claims := &jwtClaims{
			Rollno: user.Rollno,
			Role:   role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiry.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		var response Token
		response.Jwt = tokenString
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(response)
		addLog(fmt.Sprintf("%s%d", "Log in by: ", user.Rollno))
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(fmt.Sprintln("Only Post Method allowed")))
	}
}
func signup(res http.ResponseWriter, req *http.Request) {
	var user signupRoute
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		if user.Password == "" || user.Rollno == 0 || user.Username == "" {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}

		password := []byte(user.Password)
		hash, err := bcrypt.GenerateFromPassword(password, 16)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader((http.StatusInternalServerError))
			res.Write([]byte(fmt.Sprintf("Error:%d Internal Server Error", http.StatusInternalServerError)))
			return
		}
		user.Password = string(hash)
		err = AddStudent(user)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader((http.StatusInternalServerError))
			res.Write([]byte(fmt.Sprintln("Already Registered")))
			return
		}
		res.WriteHeader(http.StatusAccepted)
		addLog(fmt.Sprintf("%s%d%s%s", "New account created for roll number: ", user.Rollno, " and username: ", user.Username))
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(("Only Post Method allowed")))
	}
}
func changePass(res http.ResponseWriter, req *http.Request) {
	var user changePassRoute
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		old, _ := bcrypt.GenerateFromPassword([]byte(user.Old), 16)
		ne, _ := bcrypt.GenerateFromPassword([]byte(user.New), 16)
		err = changePassword(user.Rollno, string(old), string(ne))
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusAccepted)
		addLog(fmt.Sprintf("%s%d", "Password changed by: ", user.Rollno))
		return

	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(fmt.Sprintln("Only Post Method allowed")))
	}
}
func changeRole(res http.ResponseWriter, req *http.Request) {
	var user changeRoleRoute
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&user)
		if err != nil {
			res.WriteHeader((http.StatusBadRequest))
			res.Write([]byte(fmt.Sprintf("Error:%d Bad Request", http.StatusBadRequest)))
			return
		}
		if user.Role != 0 && user.Role != 1 && user.Role != 2 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &jwtClaims{}
		tkn, err := jwt.ParseWithClaims(user.Jwt, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims.Role != 3 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = changeRoleDB(user)
		if err != nil {
			res.Write([]byte("Some error occured please try again later"))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusAccepted)
		addLog(fmt.Sprintf("%s%d%s%d%s%d", "Role of user ", user.Rollno, " changed by: ", claims.Rollno, " to ", user.Role))
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte(fmt.Sprintln("Only Post Method allowed")))
	}
}
func changePassword(roll int, old string, new string) error {
	_, err := db.Exec("UPDATE Users SET password = ? WHERE rollno=? AND password=?", new, roll, old)
	return err
}
func changeRoleDB(user changeRoleRoute) error {
	_, err := db.Exec("UPDATE Users SET role = ? WHERE rollno=?", user.Role, user.Rollno)
	return err
}
func AddCoin(user addRoute) error {
	_, err := db.Exec("BEGIN TRANSACTION;UPDATE Users SET coins = coins + ? WHERE rollno=?;UPDATE Users SET participation=participation+1 WHERE rollno=?;END TRANSACTION;", user.Coin, user.Rollno, user.Rollno)
	return err
}
func SendCoin(user sendRoute, senderRoll int, isSameBatch bool) {
	left_after_tax := 1.0
	if isSameBatch {
		left_after_tax = 0.98 * float64(user.Coins)
	} else {
		left_after_tax = 0.67 * float64(user.Coins)
	}
	_, err = db.Exec("END TRANSACTION")
	_, err = db.Exec("BEGIN TRANSACTION; UPDATE Users SET coins=coins-? WHERE rollno=?; UPDATE Users SET coins=coins+? WHERE rollno=?;COMMIT;", user.Coins, senderRoll, left_after_tax, user.Reciever)
	if err != nil {
		fmt.Println(err)
	}
}
func GetRole(rollno int) int {
	searchQuery := "SELECT role FROM Users WHERE rollno=?"
	row := db.QueryRow(searchQuery, rollno)
	var role int
	err = row.Scan(&role)
	if err != nil {
		fmt.Println(err)
		return -1
	} else {
		return role
	}
}
func GetBalance(rollno int) int {
	searchQuery := "SELECT coins FROM Users WHERE rollno=?"
	tmp, err := db.Prepare(searchQuery)
	if err != nil {
		fmt.Println(err)
	}
	var coins int
	err = tmp.QueryRow(rollno).Scan(&coins)
	if err != nil {
		return 0
	}
	return coins
}
func SearchStudent(user loginRoute) int {
	searchQuery := "SELECT password,role FROM Users WHERE rollno=?"
	row := db.QueryRow(searchQuery, user.Rollno)
	var role int
	var hash string
	err = row.Scan(&hash, &role)
	if err != nil {
		fmt.Println(err)
		return -1
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
		if err != nil {
			return -1
		}
		return role
	}
}
func AddStudent(user signupRoute) error {

	addQuery := "INSERT OR IGNORE INTO Users(rollno, password,username,coins,role,participation) VALUES (?, ?,?,?,?,?)"

	_, err = db.Exec(addQuery, user.Rollno, user.Password, user.Username, 0, 0, 0)

	if err != nil {
		fmt.Println(err)
	}
	return err
}
func addLog(operation string) {
	addQuery := "INSERT INTO logs (operation) VALUES (?)"
	_, err = db.Exec(addQuery, operation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func startDB() error {
	createTableQuery := "CREATE TABLE IF NOT EXISTS Users (rollno INTEGER, password TEXT,username TEXT,coins FLOAT DEFAULT 0,role INTEGER DEFAULT 0,participation INTEGER DEFAULT 0,CHECK(coins>=0),CHECK(coins<=1000000),UNIQUE(rollno),UNIQUE(username))"

	_, err := db.Exec(createTableQuery)
	if err != nil {
		fmt.Println(err)
		return err
	}
	createTableQuery = "CREATE TABLE IF NOT EXISTS logs (timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,operation TEXT)"

	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

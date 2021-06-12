package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t User
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	password := []byte(t.Password)
	hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	exists := search(t.Rollno, string(hash))
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expiry := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Rollno: t.Rollno,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	var response Token
	response.Jwt = tokenString
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func signupHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t User
	_ = decoder.Decode(&t)
	password := []byte(t.Password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	t.Password = string(hash)
	addStudent(t)
	w.WriteHeader(http.StatusAccepted)
}
func secrethandler(w http.ResponseWriter, r *http.Request) {
	var tokens Token
	decoder := json.NewDecoder(r.Body)
	var t User
	_ = decoder.Decode(&t)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokens.Jwt, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Printf("Secret Page %x", claims.Rollno)
}

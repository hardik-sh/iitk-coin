package main

import "github.com/dgrijalva/jwt-go"

var jwtKey = []byte("iitk-coin")

type jwtClaims struct {
	Rollno int `json:"rollno"`
	Role   int `json:"role"`
	jwt.StandardClaims
}
type loginRoute struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}
type signupRoute struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
	Username string `json:"username"`
}
type Token struct {
	Jwt string `json:"jwt"`
}
type Coins struct {
	Balance int `json:"balance"`
}
type addRoute struct {
	Jwt    string `json:"jwt"`
	Rollno int    `json:"rollno"`
	Coin   int    `json:"coins"`
}
type changeRoleRoute struct {
	Jwt    string `json:"jwt"`
	Rollno int    `json:"rollno"`
	Role   int    `json:"role"`
}
type sendRoute struct {
	SenderJWT string `json:"jwt"`
	Coins     int    `json:"coins"`
	Reciever  int    `json:"reciever"`
}
type changePassRoute struct {
	Rollno int    `json:"rollno"`
	Old    string `json:"old"`
	New    string `json:"new"`
}

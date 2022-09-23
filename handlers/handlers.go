package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// Define a variable of type byte for secret key
var jwtkey = []byte("secret_key")

// Create a map to store users name and their passwords
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// create struct for credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// create  claim struct for jwt claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Login function receive request in r variable and write response in w
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials                         // Create variable of type credentials struct
	err := json.NewDecoder(r.Body).Decode(&credentials) // it receive credentials in json format and decode it
	if err != nil {                                     // if err occur during decoding of json credentials then it execute
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users[credentials.Username]
	if !ok || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// now set expirationtime 5mint for token
	expirationTime := time.Now().Add(time.Minute * 5)
	//set claims for login user
	claims := &Claims{
		Username: credentials.Username, // current login user name
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //user token allowed time
		},
	}
	//create token and pass above claim in it
	// SigningMethodHS256 is encryption algo, you can choose algo according to your choice
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey) //set singed string for a token
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//create a cookie as respone write and pass token, tokenstring and expirationtime in it for login user
	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
}

func Home(w http.ResponseWriter, r *http.Request) {
	//will recive token after click on login button
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get values from cookie and assign it to tokenstr variable
	tokenStr := cookie.Value
	claims := &Claims{}                               //create claims varialble and pass reference of Claims struct
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, //parse the token
		func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid { //if token is not valid
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	// create refresh token for login user if current token time remain less than 30 sec
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtkey, nil
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
	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	expirationTime := time.Now().Add(time.Minute * 5)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,
		&http.Cookie{
			Name:    "refresh_token",
			Value:   tokenString,
			Expires: expirationTime,
		})
}

package handlers

import (
	"GolangJWTService/src/auth"
	"GolangJWTService/src/errors"
	"GolangJWTService/src/models"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

// GetToken serving token handler
func GetToken(w http.ResponseWriter, r *http.Request) {

	NotPost(w, r)
	SetContentJson(w)

	if r.Header["X-Api-Key"] != nil {

		if auth.CheckPassword(r.Header["X-Api-Key"][0], []byte(models.ApiKey{}.GetXApiKey().XApiKey)) {

			token, err := CreateToken()
			errors.CheckErr(err)

			log.Default().Print("success token")

			m := make(map[string]string)
			m["token"] = token
			m["status"] = "success"

			jsonS, _ := json.Marshal(m)

			_, err = w.Write(jsonS)
			errors.CheckErr(err)

		} else {
			SendUnAuthWrite(w)
		}
	} else {
		SendUnAuthWrite(w)
	}
}

// CreateToken this function generates tokens
func CreateToken() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString([]byte(models.ApiKey{}.GetXApiKey().SecretKey))
	errors.CheckErr(err)

	return tokenStr, nil
}

package handlers

import (
	"GolangJWTService/src/errors"
	"encoding/json"
	"log"
	"net/http"
)

// Index serving index handler
func Index(w http.ResponseWriter, r *http.Request) {

	NotPost(w, r)

	SetContentJson(w)

	log.Default().Print("success access home page")

	m := make(map[string]string)
	m["status"] = "success"
	m["message"] = "home page"

	jsonM, _ := json.Marshal(m)

	var _, err = w.Write(jsonM)

	errors.CheckErr(err)
}

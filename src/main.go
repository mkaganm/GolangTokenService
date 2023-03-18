package main

import (
	"GolangJWTService/src/errors"
	handlers "GolangJWTService/src/handlers"
	"log"
	"net/http"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc(
		"/token",
		handlers.GetToken,
	)
	r.Handle(
		"/",
		handlers.ValidateToken(handlers.Index),
	)

	log.Default().Print("Server localhost:8081 is started")

	defer func(addr string, handler http.Handler) {
		err := http.ListenAndServe(addr, handler)
		errors.CheckErr(err)
	}(":8081", r)
}

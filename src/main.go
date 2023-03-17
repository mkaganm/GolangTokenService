package main

import (
	"GolangJWTService/src/auth"
	"fmt"
)

func main() {
	apiKey := "123456"

	fmt.Println(auth.PasswordToHash(apiKey))
}
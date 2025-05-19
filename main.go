package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("7"), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	fmt.Println(string(hashedBytes))

}

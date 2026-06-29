package main

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/application/session"
)

func main() {
	gameSession, err := session.NewSession()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	err = gameSession.Run()
	if err != nil {
		fmt.Print(err.Error())
	}
}

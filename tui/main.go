package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gowordle.com/display"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// no .env file is fine, env vars may be set externally
	}
	if err := display.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

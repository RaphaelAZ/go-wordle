package main

import (
	"fmt"
	"os"

	"gowordle.com/display"
)

func main() {
	if err := display.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

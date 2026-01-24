package main

import (
	"os"

	"github.com/itda-jindo/jindo/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

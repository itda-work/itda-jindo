package main

import (
	"os"

	"github.com/itda-work/itda-jindo/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

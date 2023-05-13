package main

import (
	"fmt"
	"os"

	"github.com/sourcegraph/lsif-go/cmd/lsif-go/api"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("error: %v\n", err))
		os.Exit(1)
	}
}

func mainErr() (err error) {
	return api.MainArgs(os.Args[1:])
}

package main

import (
	"fmt"
	"os"

	"github.com/elliotxx/osp/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

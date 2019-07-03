package main

import (
	"os"

	"github.com/d-kuro/approve-bot/cmd"
)

func main() {
	os.Exit(cmd.Execute(os.Stdout, os.Stderr))
}

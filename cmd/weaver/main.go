package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const modDirectory = "mods/"

func main() {
	app := &cli.App{
		Name:   "weaver",
		Usage:  "Install all mods from the `weaver.toml` file.",
		Action: Install,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

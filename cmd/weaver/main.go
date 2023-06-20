package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const modDirectory = "mods/"

func main() {
	app := &cli.App{
		Name: "weaver",
		Usage: "Minecraft server manager",
		Commands: []*cli.Command{
			{
				Name: "install",
				Aliases: []string{"i"},
				Usage: "install all mods from the toml file",
				Action: Install,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

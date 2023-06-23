package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const ConfigFileName = "weaver.toml"
const ModDirectory = "mods/"

func main() {
	app := &cli.App{
		Name:   "weaver",
		Usage:  "Minecraft Fabric server manager",
		Action: Install,
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install all mods",
				Action:  Install,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add a mod to the list",
				Action:  Add,
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "Remove a mod from the list",
				Action:  Remove,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

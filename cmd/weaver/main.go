package main

import (
	"fmt"
	"git.sr.ht/~ansipunk/weaver/pkg/cfg"
	"log"
)

func main() {
	config, configErr := cfg.ReadConfig("weaver.toml")

	if configErr != nil {
		log.Fatal(configErr)
	}

	fmt.Println(config.Loader)
	fmt.Println(config.GameVersion)
	fmt.Println(config.Mods)
}

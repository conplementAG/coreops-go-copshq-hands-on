package main

import (
	"github.com/conplementag/cops-hq/v2/pkg/hq"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	hq.GetCli().AddBaseCommand("get-config", "get-config command", "prints config values", func() {
		// 1. login to azure

		// 2. load config

		// 3. print config values
	})
}

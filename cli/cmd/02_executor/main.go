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
	hq.GetCli().AddBaseCommand("get-version", "get version command", "print version of az cli", func() {
		// 1. execute 'az version' command

		// 2. unmarshal json result to AzureCliVersion
	})
}

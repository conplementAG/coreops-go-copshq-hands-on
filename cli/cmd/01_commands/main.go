package main

import (
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/sirupsen/logrus"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	// 1. add hello world command
	hq.GetCli().AddBaseCommand("hello", "prints greeting", "prints greeting to the stdout", func() {
		logrus.Info("Hello, World!")
	})

	// 2. add welcome subcommand to hello command

	// 3. add parameter so hello command

	// 4. add bye subcommand to hello command
}

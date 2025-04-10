package main

import (
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	// 1. add hello world command
	helloCommand := hq.GetCli().AddBaseCommand("hello", "prints greeting", "prints greeting to the stdout", func() {
		logrus.Info("Hello, World!")
	})

	// 2. add welcome subcommand to hello command
	helloCommand.AddCommand("welcome", "prints welcome", "prints welcome to stdout", func() {
		logrus.Info("Welcome to the workshop!")
	})

	// 3. add parameter so hello command
	helloCommand.AddPersistentParameterString("name", "Conplement", false, "", "name of the person to greet")

	// 4. add bye subcommand to hello command
	helloCommand.AddCommand("bye", "prints bye", "prints bye to stdout", func() {
		logrus.Infof("Bye, %s!", viper.GetString("name"))
	})
}

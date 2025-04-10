package main

import (
	"github.com/conplementag/cops-hq/v2/pkg/hq"
	"github.com/conplementag/cops-hq/v2/pkg/recipes/azure_login"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	hq := hq.NewQuiet("workshop-cli", "0.0.1", "workshop-cli.log")

	createCommands(hq)

	hq.Run()
}

func createCommands(hq hq.HQ) {
	hq.GetCli().AddBaseCommand("get-config", "get-config command", "prints config values", func() {
		login := azure_login.New(hq.GetExecutor())
		err := login.Login()
		if err != nil {
			panic(err)
		}

		err = hq.LoadConfigFile("./config/config.yaml")
		if err != nil {
			panic(err)
		}
		logrus.Info(viper.GetString("hello"))
		logrus.Info(viper.GetString("example_secret"))
	})
}

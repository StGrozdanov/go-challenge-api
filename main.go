package main

import (
	log "github.com/sirupsen/logrus"
	"outdoorsy-api/config"
	"outdoorsy-api/server"
	"outdoorsy-api/utils"
)

func init() {
	app, err := config.Init()
	if err != nil {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on config initialization")
		return
	}
	if app.AppEnv == "LOC" {
		utils.PrettyPrint(app)
	}
}

func main() {
	server.Run()
}

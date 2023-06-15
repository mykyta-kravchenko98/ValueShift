package main

import (
	"log"
	"os"
	"valueShift/cmd/commands"
	"valueShift/internal/configs"
	"valueShift/internal/db"
	"valueShift/internal/db/repositories"
	"valueShift/internal/services"
)

const (
	DBName = "valueshift"
)

func main() {
	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	//Load configuration
	conf, confErr := configs.LoadConfigs(env)
	if confErr != nil {
		log.Fatal(confErr)
	}

	// Setup : DB
	dbManager, dErr := db.NewMongoManager(DBName, conf.MongoDB.URL)
	if dErr != nil {
		log.Fatal(dErr)
	}

	rep := repositories.NewCurrencySnapshotDataService(dbManager.Database())
	svc := services.NewCurrencySnapshotDataService(rep)

	//Load Cobra
	commands.InitServices(svc)
	commands.Execute()
}

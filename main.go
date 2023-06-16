package main

import (
	"log"
	"os"
	"time"
	"valueShift/cmd/commands"
	_ "valueShift/docs"
	"valueShift/internal/configs"
	"valueShift/internal/db"
	"valueShift/internal/db/repositories"
	"valueShift/internal/models"
	"valueShift/internal/server"
	"valueShift/internal/services"
)

const (
	DBName      = "valueshift"
	ServiceName = "valueshift-api"
)

// Passed while building from  make file
var version string = "1.0.0"

// @title           API for converting currencies
// @version         1.0
// @description     API that provide abuility to converting currencies

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	upTime := time.Now()
	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	// Metadata of the service
	serviceInfo := &models.ServiceInfo{
		Name:        ServiceName,
		UpTime:      upTime,
		Environment: env,
		Version:     version,
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

	//Server Init
	server.Init(serviceInfo, dbManager)
}

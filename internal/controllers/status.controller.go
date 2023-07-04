package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mykyta-kravchenko98/ValueShift/internal/db"
	"github.com/mykyta-kravchenko98/ValueShift/internal/models"

	"github.com/gin-gonic/gin"
)

type ServiceStatus string

type StatusResponce struct {
	Status      ServiceStatus
	ServiceName string
	UpTime      time.Time
	Environment string
	Version     string
}

type StatusController struct {
	svcInfo   *models.ServiceInfo
	dbManager db.MongoManager
}

const (
	UP   ServiceStatus = "OK"
	DOWN ServiceStatus = "DOWN"
)

var (
	UnableConnectToDb = errors.New("nable to connect to DB")
)

// New instance init
func NewStatusController(s *models.ServiceInfo, dbManager db.MongoManager) *StatusController {
	return &StatusController{
		svcInfo:   s,
		dbManager: dbManager,
	}
}

// Checks the health of all the dependencies of the service to ensure complete serviceability
func (sc *StatusController) CheckStatus(ctx *gin.Context) {
	var status ServiceStatus
	var statusCode int

	if err := sc.dbManager.Ping(); err != nil {
		log.Print(UnableConnectToDb)
		status = DOWN
		statusCode = http.StatusFailedDependency
	} else {
		status = UP
		statusCode = http.StatusOK
	}

	statusResponce := StatusResponce{
		Status:      status,
		ServiceName: sc.svcInfo.Name,
		UpTime:      sc.svcInfo.UpTime,
		Environment: sc.svcInfo.Environment,
		Version:     sc.svcInfo.Version,
	}

	ctx.JSON(statusCode, statusResponce)
}

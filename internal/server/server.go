package server

import (
	"sync"

	"github.com/mykyta-kravchenko98/ValueShift/internal/configs"
	"github.com/mykyta-kravchenko98/ValueShift/internal/controllers"
	"github.com/mykyta-kravchenko98/ValueShift/internal/db"
	"github.com/mykyta-kravchenko98/ValueShift/internal/db/repositories"
	"github.com/mykyta-kravchenko98/ValueShift/internal/models"
	"github.com/mykyta-kravchenko98/ValueShift/internal/services"
	"github.com/mykyta-kravchenko98/ValueShift/util"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var runOnce sync.Once

func Init(serviceInfo *models.ServiceInfo, manager db.MongoManager) {
	config := configs.GetConfig()
	port := config.Server.Port
	runOnce.Do(func() {
		r := WebRouter(serviceInfo, manager)
		r.Run(":" + port)
	})
}

func WebRouter(serviceInfo *models.ServiceInfo, manager db.MongoManager) (router *gin.Engine) {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(serviceInfo.Environment) {
		ginMode = gin.DebugMode
		gin.ForceConsoleColor()
	}

	gin.SetMode(ginMode)

	//Middleware
	router = gin.Default()
	pprof.Register(router)

	//Routes Status Controller
	status := controllers.NewStatusController(serviceInfo, manager)
	router.GET("/status", status.CheckStatus)

	// Dependencies for controllers
	database := manager.Database()
	currencySnapshotDataService := repositories.NewCurrencySnapshotDataService(database)
	currencyService := services.NewCurrencySnapshotDataService(currencySnapshotDataService)

	//Routes Convert Controller
	v1 := router.Group("/api/v1")
	{
		convertGroup := v1.Group("convert")
		{
			convert := controllers.NewConvertController(currencyService)
			convertGroup.POST("", convert.Post)
		}
	}

	//Routes - Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return
}

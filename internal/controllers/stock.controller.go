package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mykyta-kravchenko98/ValueShift/pkg/clients/grpc"
)

type stockController struct {
	yahooFinanceService grpc.YahooFinanceService
}

func NewStockController(svc grpc.YahooFinanceService) *stockController {
	return &stockController{
		yahooFinanceService: svc,
	}
}

// GetAll  godoc
// @Summary      Convert currency
// @Description  Used to convert one currency to another
// @ID GetAll
// @Param data body ConvertDto true "The body to do converting"
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /stock/all [get]
func (stockController *stockController) GetAll(ctx *gin.Context) {
	result, err := stockController.yahooFinanceService.GetAllValidStocks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, result)
}

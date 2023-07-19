package controllers

import (
	"net/http"

	"github.com/mykyta-kravchenko98/ValueShift/internal/services"

	"github.com/gin-gonic/gin"
)

type currencyController struct {
	currencyService services.CurrencyService
}

func NewCurrencyController(svc services.CurrencyService) *currencyController {
	return &currencyController{
		currencyService: svc,
	}
}

// Get  godoc
// @Summary      List of avalible cyrrencies
// @Description  Used to get list of avalible cyrrencies
// @ID Post
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /currency/list [get]
func (currencyController *currencyController) GetList(ctx *gin.Context) {
	result, err := currencyController.currencyService.GetCurrencyList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, result)
}

package controllers

import (
	"errors"
	"net/http"

	"github.com/mykyta-kravchenko98/ValueShift/internal/services"

	"github.com/gin-gonic/gin"
)

type ConvertController struct {
	currencyService services.CurrencyService
}

func NewConvertController(svc services.CurrencyService) *ConvertController {
	return &ConvertController{
		currencyService: svc,
	}
}

// Dto for converting requests
// ConvertDto  godoc
type ConvertDto struct {
	InputCurrencyLable  string  `json:"input_currency_lable" example:"USD"`
	OutputCurrencyLable string  `json:"output_currency_lable" example:"EUR"`
	Value               float64 `json:"value" example:"3000"`
} // @name ConvertDto

var (
	BadRequest = errors.New("")
)

// Post  godoc
// @Summary      Convert currency
// @Description  Used to convert one currency to another
// @ID Post
// @Param data body ConvertDto true "The body to do converting"
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /convert/ [post]
func (convertController *ConvertController) Post(ctx *gin.Context) {
	requestDto := ConvertDto{}

	if err := ctx.BindJSON(&requestDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	if err := validateConvertDto(requestDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	result, err := convertController.currencyService.Converting(requestDto.InputCurrencyLable, requestDto.OutputCurrencyLable, requestDto.Value)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, result)
}

func validateConvertDto(object ConvertDto) error {
	if object.InputCurrencyLable == "" {
		return errors.New("input_currency_lable can`t be empty,")
	}

	if object.OutputCurrencyLable == "" {
		return errors.New("output_currency_lable can`t be empty.")
	}

	if object.Value <= 0 {
		return errors.New("value can`t be equal or less that zero.")
	}

	return nil
}

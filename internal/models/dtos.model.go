package models

// GetCurrencyListResponseDto  godoc
type GetCurrencyListResponseDto struct {
	Lables []string `json:"labels" example:"[USD, EUR, UAH]"`
} // @name GetCurrencyListResponseDto

// Dto for converting requests
// ConvertDto  godoc
type ConvertDto struct {
	InputCurrencyLable  string  `json:"input_currency_lable" example:"USD"`
	OutputCurrencyLable string  `json:"output_currency_lable" example:"EUR"`
	Value               float64 `json:"value" example:"3000"`
} // @name ConvertDto

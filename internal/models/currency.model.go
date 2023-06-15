package models

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrencySnapshot struct {
	Id              primitive.ObjectID     `bson:"_id,omitempty"`
	Lable           string                 `bson:"lable" json:"base_code"`
	LastUpdate      string                 `bson:"last_update" json:"time_last_update_utc"`
	LastUpdateUnix  int64                  `bson:"last_update_unix" json:"time_last_update_unix"`
	NextUpdate      string                 `bson:"next_update" json:"time_next_update_utc"`
	NextUpdateUnix  int64                  `bson:"next_update_unix" json:"time_next_update_unix"`
	ConversionRates map[string]interface{} `bson:"conversion_rates" json:"conversion_rates"`
}

func (currencySnapshot *CurrencySnapshot) Converting(inputCurrencyLable, outputCurrencyLable string, value float64) (float64, error) {
	if value <= 0 {
		return -1, errors.New("Value can not be less or equel to zero.")
	}
	if inputCurrencyLable == "" || outputCurrencyLable == "" {
		return -1, errors.New("Currency lable can`t be empty")
	}

	if inputCurrencyLable == currencySnapshot.Lable {
		rate, err := currencySnapshot.parseMap(outputCurrencyLable)

		if err != nil {
			return -1, err
		}

		return value * rate, nil
	}

	if outputCurrencyLable == currencySnapshot.Lable {
		rate, err := currencySnapshot.parseMap(inputCurrencyLable)

		if err != nil {
			return -1, err
		}

		return value / rate, nil
	}

	return -1, errors.New("Wrong input or output currency.")
}

func (currencySnapshot *CurrencySnapshot) parseMap(key string) (float64, error) {
	val, ok := currencySnapshot.ConversionRates[key]
	if !ok {
		return -1, errors.New(fmt.Sprintf("Conversion Rates for %s not exist.", key))
	}
	rate, ok := val.(float64)

	if !ok {
		return -1, errors.New(fmt.Sprintf("Can not extract Conversion Rates for %s.", key))
	}

	return rate, nil
}

package models_test

import (
	"testing"

	"github.com/mykyta-kravchenko98/ValueShift/internal/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Currency_Converting_In_Right_Way(t *testing.T) {
	to := currencySnapshotBuilder()
	result, err := to.Converting("USD", "EUR", 3000)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1500.0, result)
}

func Test_Currency_Converting_Find_Input_Currency_In_ConversionRates_Collection_When_It_Possible(t *testing.T) {
	to := currencySnapshotBuilder()
	result, err := to.Converting("EUR", "USD", 3000)
	assert.Equal(t, nil, err)
	assert.Equal(t, 6000.0, result)
}

func Test_Currency_Cant_Converting_Because_Currency_Not_Exist_In_ConversionRates_Collection(t *testing.T) {
	to := currencySnapshotBuilder()

	result, err := to.Converting("USD", "UAH", 3000)

	assert.NotNil(t, err)
	assert.Equal(t, -1.0, result)
}

func Test_Currency_Value_Cant_Be_Equal_Or_Lower_Then_Zero(t *testing.T) {
	to := currencySnapshotBuilder()

	result, err := to.Converting("USD", "UAH", -100)

	assert.NotNil(t, err)
	assert.Equal(t, -1.0, result)
}

func Test_Currency_Lables_Cant_Be_Empty(t *testing.T) {
	to := currencySnapshotBuilder()

	result1, err1 := to.Converting("", "EUR", 1000)
	result2, err2 := to.Converting("USD", "", 1000)

	assert.NotNil(t, err1)
	assert.Equal(t, -1.0, result1)
	assert.NotNil(t, err2)
	assert.Equal(t, -1.0, result2)
}

func currencySnapshotBuilder() models.CurrencySnapshot {
	object := models.CurrencySnapshot{
		Id:             primitive.NewObjectID(),
		Lable:          "USD",
		LastUpdate:     "Thu, 15 Jun 2023 00:00:01 +0000",
		LastUpdateUnix: 1686787201,
		NextUpdate:     "Fri, 16 Jun 2023 00:00:01 +0000",
		NextUpdateUnix: 1686873601,
		ConversionRates: map[string]interface{}{
			"EUR": 0.5,
		},
	}

	return object
}

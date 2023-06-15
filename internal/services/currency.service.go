package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"valueShift/internal/configs"
	"valueShift/internal/db/repositories"
	"valueShift/internal/models"
	"valueShift/pkg/clients/rest"
)

type currencyService struct {
	currencySnapshotRepository repositories.CurrencySnapshotDataService
}

type CurrencyService interface {
	getCurrencySnapshot(inputCurrencyLable, outputCurrencyLable string) (models.CurrencySnapshot, error)
	Converting(inputCurrencyLable, outputCurrencyLable string, value float64) (float64, error)
}

var (
	config *configs.Config
)

func NewCurrencySnapshotDataService(svc repositories.CurrencySnapshotDataService) CurrencyService {
	image := &currencyService{
		currencySnapshotRepository: svc,
	}

	return image
}

func (currSvc *currencyService) getCurrencySnapshot(inputCurrencyLable, outputCurrencyLable string) (result models.CurrencySnapshot, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	//Searching in DB
	result, err = currSvc.currencySnapshotRepository.GetFirstExistCurrency(ctx, inputCurrencyLable, outputCurrencyLable)
	if err != nil && err != repositories.ErrNoDocuments {
		return result, err
	}
	//Result founded in db
	if err != repositories.ErrNoDocuments && !result.Id.IsZero() {
		return result, err
	}
	//Starting request proccess
	if config == nil {
		config = configs.GetConfig()
	}

	url := fmt.Sprintf("%s/%s/latest/%s", config.ExchangeApi.URL, config.ExchangeApi.ApiKey, inputCurrencyLable)

	resp, err := rest.Get(url)
	if err != nil {
		return models.CurrencySnapshot{}, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.CurrencySnapshot{}, err
	}

	//Saving in DB request result
	result, err = currSvc.currencySnapshotRepository.Create(ctx, result)

	return result, err
}

func (currSvc *currencyService) Converting(inputCurrencyLable, outputCurrencyLable string, value float64) (float64, error) {
	snapshot, err := currSvc.getCurrencySnapshot(inputCurrencyLable, outputCurrencyLable)

	if err != nil {
		return -1, err
	}

	result, err := snapshot.Converting(inputCurrencyLable, outputCurrencyLable, value)

	return result, err
}

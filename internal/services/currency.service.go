package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mykyta-kravchenko98/ValueShift/internal/configs"
	"github.com/mykyta-kravchenko98/ValueShift/internal/db/repositories"
	"github.com/mykyta-kravchenko98/ValueShift/internal/models"
	"github.com/mykyta-kravchenko98/ValueShift/pkg/clients/rest"
)

type currencyService struct {
	currencySnapshotRepository repositories.CurrencySnapshotDataService
}

type CurrencyService interface {
	Converting(inputCurrencyLable, outputCurrencyLable string, value float64) (float64, error)
	GetCurrencyList() (models.GetCurrencyListResponseDto, error)
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

	resp, err := rest.Get(url, nil)
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

func (currSvc *currencyService) GetCurrencyList() (models.GetCurrencyListResponseDto, error) {
	snapshot, err := currSvc.getCurrencySnapshot("USD", "EUR")

	if err != nil {
		return models.GetCurrencyListResponseDto{}, err
	}

	lables := make([]string, 0, len(snapshot.ConversionRates))

	for lable, _ := range snapshot.ConversionRates {
		lables = append(lables, lable)
	}

	return models.GetCurrencyListResponseDto{Lables: lables}, err
}

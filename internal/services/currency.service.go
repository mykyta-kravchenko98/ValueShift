package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"valueShift/internal/configs"
	"valueShift/internal/db/repositories"
	"valueShift/internal/models"
)

type currencyService struct {
	currencySnapshotRepository repositories.CurrencySnapshotDataService
}

type CurrencyService interface {
	GetCurrencySnapshot(lable string) (models.CurrencySnapshot, error)
}

var (
	defaultClient *http.Client
	config        *configs.Config
)

func init() {
	defaultClient = &http.Client{}
}

func NewCurrencySnapshotDataService(svc repositories.CurrencySnapshotDataService) CurrencyService {
	image := &currencyService{
		currencySnapshotRepository: svc,
	}

	return image
}

func (currSvc *currencyService) GetCurrencySnapshot(lable string) (result models.CurrencySnapshot, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	//Searching in DB
	result, err = currSvc.currencySnapshotRepository.GetByCurrency(ctx, lable)
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

	url := fmt.Sprintf("%s/%s/latest/%s", config.ExchangeApi.URL, config.ExchangeApi.ApiKey, lable)
	fmt.Println(url)

	resp, err := defaultClient.Get(url)
	if err != nil {
		return models.CurrencySnapshot{}, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.CurrencySnapshot{}, err
	}

	//Saving in DB request result
	uid, err := currSvc.currencySnapshotRepository.Create(ctx, result)

	fmt.Println(uid)

	return result, err
}

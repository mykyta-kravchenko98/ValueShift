package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"valueShift/internal/configs"
	"valueShift/internal/db"
	"valueShift/internal/db/repositories"
	"valueShift/internal/services"
	"valueShift/pkg/clients/rest"
	"valueShift/pkg/clients/rest/mock"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/assert"
)

var (
	ts        services.CurrencyService
	tr        repositories.CurrencySnapshotDataService
	dbManager db.MongoManager
)

func init() {
	rest.Client = &mock.MockRestClient{}

	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	conf, err := configs.LoadConfigs("dev")
	if err != nil {
		log.Fatal(err)
	}
	// Setup : DB
	dbManager, err = db.NewMongoManager("test", conf.MongoDB.URL)
	if err != nil {
		log.Fatal(err)
	}

	tr = repositories.NewCurrencySnapshotDataService(dbManager.Database())
	ts = services.NewCurrencySnapshotDataService(tr)
}

func Test_Converting_Process(t *testing.T) {
	err := dbManager.Database().Collection(repositories.CurrencySnapshotCollections).Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	json := `{
		"result": "success",
		"documentation": "https://www.exchangerate-api.com/docs",
		"terms_of_use": "https://www.exchangerate-api.com/terms",
		"time_last_update_unix": 1686787201,
		"time_last_update_utc": "Thu, 15 Jun 2023 00:00:01 +0000",
		"time_next_update_unix": 1686873601,
		"time_next_update_utc": "Fri, 16 Jun 2023 00:00:01 +0000",
		"base_code": "USD",
		"conversion_rates": {
			"EUR": 0.5
		}
	  }`

	mock.GetDoFunc = func(*http.Request) (*http.Response, error) {
		resp := ioutil.NopCloser(bytes.NewReader([]byte(json)))
		return &http.Response{
			StatusCode: 200,
			Body:       resp,
		}, nil
	}

	firstDbCheckCount, err := dbManager.Database().Collection(repositories.CurrencySnapshotCollections).CountDocuments(context.TODO(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, int64(0), firstDbCheckCount)

	result, err := ts.Converting("USD", "EUR", 2000)

	secondDbCheckCount, err := dbManager.Database().Collection(repositories.CurrencySnapshotCollections).CountDocuments(context.TODO(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, int64(1), secondDbCheckCount)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1000.0, result)
}

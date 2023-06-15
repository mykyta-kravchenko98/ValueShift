package db_test

import (
	"log"
	"os"
	"testing"

	"valueShift/internal/configs"
	"valueShift/internal/db"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDBMgr db.MongoManager
)

func TestMain(m *testing.M) {
	conf, err := configs.LoadConfigs("dev")
	if err != nil {
		log.Fatal(err)
	}

	d, dErr := db.NewMongoManager("test", conf.MongoDB.URL)

	if dErr != nil {
		log.Fatal(dErr)
	}
	defer d.Disconnect()
	testDBMgr = d
	//insertTestData()

	os.Exit(m.Run())
}

// func insertTestData() {
// 	database := testDBMgr.Database()
// 	dSvc := db.NewOrderDataService(database)

// 	for i := 0; i < 500; i++ {
// 		currencySnapshot := []models.CurrencySnapshot{
// 			{
// 				Name:      faker.Name(),
// 				Price:     (uint)(rand.Intn(90) + 10),
// 				Remarks:   faker.Sentence(),
// 				UpdatedAt: faker.TimeString(),
// 			},
// 			{
// 				Name:      faker.Name(),
// 				Price:     (uint)(rand.Intn(1000) + 10),
// 				Remarks:   faker.Sentence(),
// 				UpdatedAt: faker.TimeString(),
// 			},
// 		}
// 		_, err := dSvc.Create(context.TODO(), currencySnapshot)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

func TestDatabase(t *testing.T) {
	d := testDBMgr.Database()
	assert.NotNil(t, d)
	assert.IsType(t, &mongo.Database{}, d)
}

func TestPing(t *testing.T) {
	err := testDBMgr.Ping()
	assert.Nil(t, err)
}

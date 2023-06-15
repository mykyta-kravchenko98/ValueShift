package repositories

import (
	"context"
	"errors"
	"time"
	database "valueShift/internal/db"
	"valueShift/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoDocuments = errors.New("Document does not exist")
)

const (
	CurrencySnapshotCollections = "currencysnapshots"
	PageSize                    = 100
)

type CurrencySnapshotDataService interface {
	Create(ctx context.Context, currencySnapshot models.CurrencySnapshot) (*mongo.InsertOneResult, error)
	GetByDate(ctx context.Context, date string) ([]models.CurrencySnapshot, error)
	GetByCurrency(ctx context.Context, lable string) (models.CurrencySnapshot, error)
}

func NewCurrencySnapshotDataService(db database.MongoDatabase) CurrencySnapshotDataService {
	iDBSvc := &currencyRepository{
		collection: db.Collection(CurrencySnapshotCollections),
	}
	return iDBSvc
}

// currencyRepository implements CurrencySnapshotDataService
type currencyRepository struct {
	collection *mongo.Collection
}

func (currSnapDataSvc *currencyRepository) Create(ctx context.Context, currencySnap models.CurrencySnapshot) (*mongo.InsertOneResult, error) {
	if vErr := validate(currSnapDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	result, err := currSnapDataSvc.collection.InsertOne(ctx, currencySnap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (currSnapDataSvc *currencyRepository) GetByDate(ctx context.Context, date string) ([]models.CurrencySnapshot, error) {
	if vErr := validate(currSnapDataSvc.collection); vErr != nil {
		return nil, vErr
	}

	dateFrom, err := time.Parse("2023-11-15", date)
	if err != nil {
		return nil, err
	}
	dateTo := dateFrom.UTC().AddDate(0, 0, 1)

	filter := bson.D{
		{
			Key: "last_update_unix",
			Value: bson.D{
				{
					Key:   "$gte",
					Value: dateFrom.Unix(),
				},
				{
					Key:   "$lt",
					Value: dateTo.Unix(),
				},
			},
		},
	}

	cur, err := currSnapDataSvc.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var result = make([]models.CurrencySnapshot, cur.RemainingBatchLength())

	for cur.Next(ctx) {
		var currentResult models.CurrencySnapshot
		err := cur.Decode(&currentResult)
		if err != nil {
			return nil, err
		}
		result = append(result, currentResult)
	}

	return result, nil
}

func (currSnapDataSvc *currencyRepository) GetByCurrency(ctx context.Context, lable string) (models.CurrencySnapshot, error) {
	if vErr := validate(currSnapDataSvc.collection); vErr != nil {
		return models.CurrencySnapshot{}, vErr
	}

	today := time.Now()
	dateFrom := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	// dateFrom, err := time.Parse("01/02/2006", time.Now().UTC().Format("01/02/2006"))
	// if err != nil {
	// 	return models.CurrencySnapshot{}, err
	// }

	filter := bson.D{
		{
			Key: "last_update_unix",
			Value: bson.D{
				{
					Key:   "$gte",
					Value: dateFrom.Unix(),
				},
			},
		},
		{
			Key:   "lable",
			Value: lable,
		},
	}
	var result models.CurrencySnapshot

	err := currSnapDataSvc.collection.FindOne(ctx, filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return models.CurrencySnapshot{}, ErrNoDocuments
	} else if err != nil {
		return models.CurrencySnapshot{}, err
	}

	return result, nil
}

func validate(collection *mongo.Collection) error {
	if collection == nil {
		return errors.New("collection is not defined")
	}
	return nil
}

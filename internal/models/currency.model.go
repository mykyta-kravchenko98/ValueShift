package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CurrencySnapshot struct {
	Id              primitive.ObjectID     `bson:"_id,omitempty"`
	Lable           string                 `bson:"lable" json:"base_code"`
	LastUpdate      string                 `bson:"last_update" json:"time_last_update_utc"`
	LastUpdateUnix  int64                  `bson:"last_update_unix" json:"time_last_update_unix"`
	NextUpdate      string                 `bson:"next_update" json:"time_next_update_utc"`
	NextUpdateUnix  int64                  `bson:"next_update_unix" json:"time_next_update_unix"`
	ConversionRates map[string]interface{} `bson:"conversion_rates" json:"conversion_rates"`
}

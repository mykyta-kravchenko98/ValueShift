package models

import (
	pb "github.com/mykyta-kravchenko98/YahooFinanceScraperAPI/pkg/yahooScraper_v1"
)

type Stock struct {
	Symbol             string
	Name               string
	Price              float64
	Change             float64
	PercentChange      float64
	Volume             string
	AvgVolumeFor3Month string
	MarketCap          string
	PERatio            string
	CreatedAt          string
	CreatedUnix        int64
}

func (stock *Stock) ProtoToDomain(proto *pb.Stocks) {
	stock.Symbol = proto.Symbol
	stock.Name = proto.Name
	stock.Price = proto.Price
	stock.Change = proto.Change
	stock.PercentChange = proto.PercentChange
	stock.Volume = proto.Volume
	stock.AvgVolumeFor3Month = proto.AvgVolumeFor_3Month
	stock.MarketCap = proto.MarketCap
	stock.PERatio = proto.PeRatio
	stock.CreatedAt = proto.CreatedAt
	stock.CreatedUnix = proto.CreatedUnix
}

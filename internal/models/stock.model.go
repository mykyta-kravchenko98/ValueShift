package models

import (
	pb "github.com/mykyta-kravchenko98/YahooFinanceScraperAPI/pkg/yahooScraper_v1"
)

type Stock struct {
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name"`
	Price              float64 `json:"price"`
	Change             float64 `json:"change"`
	PercentChange      float64 `json:"percent_change"`
	Volume             string  `json:"volume"`
	AvgVolumeFor3Month string  `json:"avg_volume_for_3_month"`
	MarketCap          string  `json:"marketCap"`
	PERatio            string  `json:"pe_ratio"`
	CreatedAt          string  `json:"created_at"`
	CreatedUnix        int64   `json:"created_unix"`
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

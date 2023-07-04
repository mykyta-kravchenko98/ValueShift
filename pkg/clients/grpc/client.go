package grpc

import (
	"context"
	"time"

	"github.com/mykyta-kravchenko98/ValueShift/internal/models"

	yahooFinanceSVC "github.com/mykyta-kravchenko98/YahooFinanceScraperAPI/pkg/yahooScraper_v1"
	"google.golang.org/grpc"
)

var defaultRequestTimeout = time.Second * 10

type grpcService struct {
	grpcClient yahooFinanceSVC.YahooStocksServiceClient
}

type YahooFinanceService interface {
	GetAllValidStocks() ([]models.Stock, error)
}

// NewGRPCService creates a new gRPC user service connection using the specified connection string.
func NewGRPCService(connString string) (YahooFinanceService, error) {
	conn, err := grpc.Dial(connString, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &grpcService{grpcClient: yahooFinanceSVC.NewYahooStocksServiceClient(conn)}, nil
}

func (s *grpcService) GetAllValidStocks() ([]models.Stock, error) {
	req := &yahooFinanceSVC.GetAllValidStocksRequest{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()
	resp, err := s.grpcClient.GetAllValidStocks(ctx, req)

	if err != nil {
		return nil, err
	}

	stocks := make([]models.Stock, 0, resp.Count)

	for _, pbStock := range resp.GetStocks() {
		stock := models.Stock{}
		stock.ProtoToDomain(pbStock)

		stocks = append(stocks, stock)
	}

	return stocks, nil
}

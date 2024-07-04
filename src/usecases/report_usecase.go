package usecase

import (
	"context"
	"fmt"
	"time"

	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
)

type ReportUsecase interface {
	Salesreport(ctx context.Context, dto *request.SalesReportDTO) (*salesOrderResult, error)
}

type reportUsecase struct {
	orderRepository repository.OrderRepository
}

type salesOrderResult struct {
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}

func NewReportUsecase(
	orderRepository *repository.OrderRepository,
) ReportUsecase {
	return &reportUsecase{
		orderRepository: *orderRepository,
	}
}

func (u *reportUsecase) Salesreport(ctx context.Context, dto *request.SalesReportDTO) (result *salesOrderResult, err error) {
	var startDate, endDate time.Time
	if dto.StartDate == "" {
		startDate = time.Now().Add(-7 * 24 * time.Hour)
	} else {
		startDate, err = time.Parse("2006-01-02", dto.StartDate)
		if err != nil {
			return nil, err
		}
	}

	if dto.EndDate == "" {
		endDate = startDate.Add(8 * 24 * time.Hour)
		// instead of 7, we add 8 days so it will show 7 days
	} else {
		endDate, err = time.Parse("2006-01-02", dto.EndDate)
		if err != nil {
			return nil, err
		}
	}

	rangeDate := endDate.Sub(startDate).Hours() / 24
	if rangeDate > 90 {
		return nil, fmt.Errorf("date_range_maximum_is_90_days")
	}

	summary, err := u.orderRepository.CountReport(ctx, repository.OrderRepositoryFilter{
		StoreID:   &dto.StoreID,
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	if err != nil {
		return nil, err
	}

	return &salesOrderResult{
		TotalSales:   summary.TotalSales,
		TotalNett:    summary.TotalNett,
		TotalProduct: summary.TotalProduct,
	}, nil
}

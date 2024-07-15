package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/gotidy/ptr"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ReportUsecase interface {
	OrderSalesReport(ctx context.Context, dto *request.OrderSalesReportDTO) (*SalesOrderResult, error)
	ProductSalesReport(ctx context.Context, dto *request.ProductSalesRxeportDTO) (*[]ProductOrderResult, *int, error)
}

type reportUsecase struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

type SalesOrderResult struct {
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}

type ProductOrderResult struct {
	ProductKey  string
	ProductName string
	QtySales    int
	QtyStock    float64
}

func NewReportUsecase(
	orderRepository *repository.OrderRepository,
) ReportUsecase {
	return &reportUsecase{
		orderRepository: *orderRepository,
	}
}

func (u *reportUsecase) dateReportValidator(start string, end string) (startDate *time.Time, endDate *time.Time, err error) {
	if start == "" {
		startDate = ptr.Of(time.Now().Add(-7 * 24 * time.Hour))
	} else {
		startDateTime, err := time.Parse("2006-01-02", start)
		if err != nil {
			return nil, nil, err
		}
		startDate = &startDateTime
	}

	if end == "" {
		endDate = ptr.Of(startDate.Add(8 * 24 * time.Hour))
		// instead of 7, we add 8 days so it will show 7 days
	} else {
		endDateTime, err := time.Parse("2006-01-02", end)
		endDateTime = endDateTime.Add(24 * time.Hour)
		if err != nil {
			return nil, nil, err
		}
		endDate = &endDateTime
	}

	rangeDate := endDate.Sub(*startDate).Hours() / 24
	if rangeDate > 90 {
		return nil, nil, fmt.Errorf("date_range_maximum_is_90_days")
	}

	return
}

func (u *reportUsecase) OrderSalesReport(ctx context.Context, dto *request.OrderSalesReportDTO) (*SalesOrderResult, error) {
	startDate, endDate, err := u.dateReportValidator(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}

	summary, err := u.orderRepository.OrderSalesReport(ctx, repository.OrderRepositoryFilter{
		StoreID:   &dto.StoreID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}

	return &SalesOrderResult{
		TotalSales:   summary.TotalSales,
		TotalNett:    summary.TotalNett,
		TotalProduct: summary.TotalProduct,
	}, nil
}

func (u *reportUsecase) ProductSalesReport(ctx context.Context, dto *request.ProductSalesRxeportDTO) (*[]ProductOrderResult, *int, error) {
	startDate, endDate, err := u.dateReportValidator(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, nil, err
	}

	page, perPage, err := utils.ParsePageAndPerPage(dto.Page, dto.PerPage)
	if err != nil {
		return nil, nil, err
	}

	summary, count, err := u.orderRepository.ProductSalesReport(ctx, page, perPage, repository.OrderRepositoryFilter{
		StoreID:   &dto.StoreID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, nil, err
	}

	var productOrderResults []ProductOrderResult
	for _, sum := range *summary {
		productOrderResults = append(productOrderResults, ProductOrderResult{
			ProductKey:  sum.ProductKey,
			ProductName: sum.ProductName,
			QtySales:    sum.QtySales,
			QtyStock:    sum.QtyStock,
		})
	}

	return &productOrderResults, count, nil
}

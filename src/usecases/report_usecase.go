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
	Salesreport(ctx context.Context, dto *request.SalesReportDTO) (*salesOrderResult, error)
}

type reportUsecase struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

type salesOrderResult struct {
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}

type productOrderResult struct {
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}

func NewReportUsecase(
	orderRepository *repository.OrderRepository,
	productRepository *repository.ProductRepository,
) ReportUsecase {
	return &reportUsecase{
		orderRepository:   *orderRepository,
		productRepository: *productRepository,
	}
}

func (u *reportUsecase) dateReportValidator(start string, end string) (startDate *time.Time, endDate *time.Time, err error) {

	// var startDate, endDate time.Time
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

func (u *reportUsecase) OrderSalesReport(ctx context.Context, dto *request.OrderSalesReportDTO) (result *salesOrderResult, err error) {
	startDate, endDate, err := u.dateReportValidator(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}

	summary, _ := u.orderRepository.OrderSalesReport(ctx, repository.OrderRepositoryFilter{
		StoreID:   &dto.StoreID,
		StartDate: startDate,
		EndDate:   endDate,
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

func (u *reportUsecase) ProductSalesReport(ctx context.Context, dto *request.ProductSalesRxeportDTO) (*[]productOrderResult, error) {
	startDate, endDate, err := u.dateReportValidator(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}

	page, perPage, err := utils.ParsePageAndPerPage(dto.Page, dto.PerPage)
	if err != nil {
		return nil, err
	}

	summary, _ := u.productRepository.ProductSalesReport(ctx, page, perPage, repository.ProductRepositoryFilter{
		StoreID:   &dto.StoreID,
		StartDate: startDate,
		EndDate:   endDate,
	})

	fmt.Println(startDate, endDate)
	fmt.Println(summary)

	return nil, nil
}

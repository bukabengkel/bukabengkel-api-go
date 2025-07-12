package admin_usecase

import (
	"context"
	"time"

	"github.com/gotidy/ptr"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	admin_request "github.com/peang/bukabengkel-api-go/src/transport/request/admin"
	admin_response "github.com/peang/bukabengkel-api-go/src/transport/response/admin"
)

type DashboardUsecase interface {
	Dashboard(ctx context.Context, dto *admin_request.AdminDashboardDTO) (*admin_response.AdminDashboardResponse, error)
}

type dashboardUsecase struct {
	orderRepository repository.OrderRepository
	productRepository repository.ProductRepository
	storeRepository repository.StoreRepository
	userRepository repository.UserStoreAggregateRepository
}

func NewDashboardUsecase(
	orderRepository *repository.OrderRepository,
	productRepository *repository.ProductRepository,
	storeRepository *repository.StoreRepository,
	userRepository *repository.UserStoreAggregateRepository,
) DashboardUsecase {
	return &dashboardUsecase{
		orderRepository: *orderRepository,
		productRepository: *productRepository,
		storeRepository: *storeRepository,
		userRepository: *userRepository,
	}
}

func (u *dashboardUsecase) Dashboard(ctx context.Context, dto *admin_request.AdminDashboardDTO) (*admin_response.AdminDashboardResponse, error) {
	now := time.Now()
	month := now.Month()
	year := now.Year()

	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

	lastMonthStart := monthStart.AddDate(0, -1, 0)
	lastMonthEnd := lastMonthStart.AddDate(0, 1, 0).Add(-time.Second)

	totalOrderThisMonth, err := u.orderRepository.Count(ctx, repository.OrderRepositoryFilter{
		StartDate: ptr.Of(monthStart),
		EndDate: ptr.Of(monthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalOrderLastMonth, err := u.orderRepository.Count(ctx, repository.OrderRepositoryFilter{
		StartDate: ptr.Of(lastMonthStart),
		EndDate: ptr.Of(lastMonthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalProductThisMonth, err := u.productRepository.Count(ctx, repository.ProductRepositoryFilter{
		StartDate: ptr.Of(monthStart),
		EndDate: ptr.Of(monthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalProductLastMonth, err := u.productRepository.Count(ctx, repository.ProductRepositoryFilter{
		StartDate: ptr.Of(lastMonthStart),
		EndDate: ptr.Of(lastMonthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalStoreThisMonth, err := u.storeRepository.Count(ctx, repository.StoreRepositoryFilter{
		StartDate: ptr.Of(monthStart),
		EndDate: ptr.Of(monthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalStoreLastMonth, err := u.storeRepository.Count(ctx, repository.StoreRepositoryFilter{
		StartDate: ptr.Of(lastMonthStart),
		EndDate: ptr.Of(lastMonthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalOrderAmountThisMonth, err := u.orderRepository.OrderSalesReport(ctx, repository.OrderRepositoryFilter{
		StartDate: ptr.Of(monthStart),
		EndDate: ptr.Of(monthEnd),
	})
	if err != nil {
		return nil, err
	}

	totalOrderAmountLastMonth, err := u.orderRepository.OrderSalesReport(ctx, repository.OrderRepositoryFilter{
		StartDate: ptr.Of(lastMonthStart),
		EndDate: ptr.Of(lastMonthEnd),
	})
	if err != nil {
		return nil, err
	}

	return &admin_response.AdminDashboardResponse{
		TotalOrderThisMonth: totalOrderThisMonth,
		TotalOrderLastMonth: totalOrderLastMonth,
		TotalProductThisMonth: totalProductThisMonth,
		TotalProductLastMonth: totalProductLastMonth,
		TotalStoreThisMonth: totalStoreThisMonth,
		TotalStoreLastMonth: totalStoreLastMonth,
		TotalOrderAmountThisMonth: totalOrderAmountThisMonth.TotalSales,
		TotalOrderAmountLastMonth: totalOrderAmountLastMonth.TotalSales,
	}, nil
}
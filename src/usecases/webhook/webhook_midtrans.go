package usecase

import (
	"context"
	"time"

	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
)

func (u *webhookUsecase) MidtransWebhook(ctx context.Context, dto request.MidtransWebhookDTO) {
	filter := repository.OrderDistributorRepositoryFilter{
		Key: &dto.OrderID,
	}

	orderDistributor, err := u.orderDistributorRepo.FindOne(ctx, filter)
	if err != nil {
		u.logger.Error(err)
	}

	// userStore, err := u.userStoreRepo.FindOne(ctx, repository.UserStoreAggregateRepositoryFilter{
	// 	UserID:  orderDistributor.Store.ID,
	// 	StoreID: orderDistributor.Store.ID,
	// })
	// if err != nil {
	// 	u.logger.Error(err)
	// }

	switch dto.TransactionStatus {
	case "pending":
		err = u.pendingPayment(ctx, orderDistributor, dto)
		if err != nil {
			u.logger.Error(err)
		}
	case "settlement":
		err = u.settlementPayment(ctx, orderDistributor, dto)
		if err != nil {
			u.logger.Error(err)
		}
	case "expire":
		err = u.expirePayment(ctx, orderDistributor, dto)
		if err != nil {
			u.logger.Error(err)
		}
	}
}

func (u *webhookUsecase) pendingPayment(ctx context.Context, orderDistributor *models.OrderDistributor, dto request.MidtransWebhookDTO) error {
	u.logger.Info("Processing PENDING Payment for Order Distributor", "order_distributor_id", orderDistributor.Key.String())

	if dto.ExpiredAt != nil {
		expiredAt, err := time.Parse(time.RFC3339, *dto.ExpiredAt)
		if err != nil {
			expiredAt, err = time.Parse("2006-01-02 15:04:05", *dto.ExpiredAt)
			if err != nil {
				return err
			}
		}
		orderDistributor.ExpiredAt = &expiredAt
	}

	orderDistributor.TransactionLogs = append(orderDistributor.TransactionLogs, models.OrderDistributorTransactionLog{
		Status:    string(models.OrderDistributorStatusWaitingForPaymentResponse),
		Timestamp: time.Now(),
		Remarks:   "Payment Type: " + dto.PaymentType,
	})
	orderDistributor.Status = models.OrderDistributorStatusWaitingForPaymentResponse
	orderDistributor.UpdatedAt = time.Now()

	err := u.orderDistributorRepo.UpdateOrderDistributor(ctx, orderDistributor)
	if err != nil {
		return err
	}

	return nil
}

func (u *webhookUsecase) settlementPayment(ctx context.Context, orderDistributor *models.OrderDistributor, dto request.MidtransWebhookDTO) error {
	u.logger.Info("Processing SETTLEMENT Payment for Order Distributor", "order_distributor_id", orderDistributor.Key.String())

	err := u.distributorService.Checkout(ctx, orderDistributor)
	if err != nil {
		orderDistributor.TransactionLogs = append(orderDistributor.TransactionLogs, models.OrderDistributorTransactionLog{
			Status:    string(models.OrderDistributorStatusRequestDistributorFailed),
			Timestamp: time.Now(),
			Remarks:   "Request Distributor Failed: " + err.Error(),
		})
		orderDistributor.Status = models.OrderDistributorStatusRequestDistributorFailed
		orderDistributor.UpdatedAt = time.Now()
	} else {
		orderDistributor.TransactionLogs = append(orderDistributor.TransactionLogs, models.OrderDistributorTransactionLog{
			Status:    string(models.OrderDistributorStatusWaitingForShipment),
			Timestamp: time.Now(),
			Remarks:   "Payment Type: " + dto.PaymentType,
		})
		orderDistributor.Status = models.OrderDistributorStatusWaitingForShipment
		orderDistributor.UpdatedAt = time.Now()
		orderDistributor.PaidAt = ptr.Of(time.Now())

		u.distributorService.Checkout(ctx, orderDistributor)
	}

	err = u.orderDistributorRepo.UpdateOrderDistributor(ctx, orderDistributor)
	if err != nil {
		return err
	}
	return nil
}

func (u *webhookUsecase) expirePayment(ctx context.Context, orderDistributor *models.OrderDistributor, dto request.MidtransWebhookDTO) error {
	u.logger.Info("Processing EXPIRED Payment for Order Distributor", "order_distributor_id", orderDistributor.Key.String())

	orderDistributor.TransactionLogs = append(orderDistributor.TransactionLogs, models.OrderDistributorTransactionLog{
		Status:    string(models.OrderDistributorStatusPaymentExpired),
		Timestamp: time.Now(),
		Remarks:   "Payment Type: " + dto.PaymentType,
	})
	return nil
}
